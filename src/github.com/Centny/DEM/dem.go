//Author:Centny
//Package DEM provide the testing sql driver.
package DEM

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

//register one drive to system by name.
func Register(n string, ev DbEv) {
	sql.Register(n, &STDriver{N: n, Ev: ev, DBS: map[*sql.DB]int{}})
}

////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////// Driver ///////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////

//the database event interface for all callback
type DbEv interface {
	//
	OnOpen(dsn string) (*sql.DB, error)
	//
	OnBegin(c *STConn) error
	OnPrepare(c *STConn, query string) error
	OnNumInput(c *STConn, query string, stm *sql.Stmt) int
	OnClose(c *STConn) error
	//
	OnTxCommit(tx *STTx) error
	OnTxRollback(tx *STTx) error
	//
	OnStmQuery(stm *STStmt, args []driver.Value) error
	OnStmExec(stm *STStmt, args []driver.Value) error
	OnStmClose(stm *STStmt) error
	//
	OnResLIId(res *STResult) error
	OnResARow(res *STResult) error
	//
	OnRowNext(row *STRows) error
	IsEmpty(row *STRows) bool
	OnRowClose(row *STRows) error
}

var C_Stack map[*STConn]string = map[*STConn]string{}
var T_Stack map[*STTx]string = map[*STTx]string{}
var S_Stack map[*STStmt]string = map[*STStmt]string{}
var R_Stack map[*STRows]string = map[*STRows]string{}
var LAST *sql.DB = nil

func CallStack() string {
	buf := make([]byte, 10240)
	l := runtime.Stack(buf, false)
	return string(buf[:l])
}

////////////////////////////////////////////////////////////////////////////////////
type STDriver struct {
	N   string //driver name.
	Ev  DbEv   //database evernt
	DBS map[*sql.DB]int
}

func (d *STDriver) Open(dsn string) (driver.Conn, error) {
	con, err := d.Ev.OnOpen(dsn)
	if err != nil {
		return nil, err
	}
	d.DBS[con] = 1
	LAST = con
	c := &STConn{
		Db: con,
		Dr: d,
		Ev: d.Ev,
		lc: sync.RWMutex{},
	}
	C_Stack[c] = CallStack()
	return c, err
}

type STConn struct {
	Db *sql.DB
	Dr *STDriver
	Ev DbEv //database evernt
	tx *STTx
	lc sync.RWMutex
}

func (c *STConn) Begin() (driver.Tx, error) {
	c.lc.Lock()
	defer c.lc.Unlock()
	if c.tx != nil {
		return nil, errors.New("already starting transaction")
	}
	if e := c.Ev.OnBegin(c); e != nil {
		return nil, e
	}
	tx, err := c.Db.Begin()
	if err != nil {
		return nil, err
	}
	c.tx = &STTx{
		Tx:   tx,
		Conn: c,
		Ev:   c.Ev,
	}
	T_Stack[c.tx] = CallStack()
	return c.tx, nil
}
func (c *STConn) TxDone() {
	c.lc.Lock()
	defer c.lc.Unlock()
	if _, ok := T_Stack[c.tx]; ok {
		delete(T_Stack, c.tx)
	} else {
		fmt.Println("done the Tx not found in stack")
	}
	c.tx = nil
}

func (c *STConn) Prepare(query string) (driver.Stmt, error) {
	c.lc.Lock()
	defer c.lc.Unlock()
	if e := c.Ev.OnPrepare(c, query); e != nil {
		return nil, e
	}
	var stm *sql.Stmt
	var err error
	if c.tx == nil {
		stm, err = c.Db.Prepare(query)
	} else {
		stm, err = c.tx.Tx.Prepare(query)
	}
	if err != nil {
		return nil, err
	}
	s := &STStmt{
		Q:    query,
		Conn: c,
		Stmt: stm,
		Ev:   c.Ev,
		Num:  c.Ev.OnNumInput(c, query, stm),
	}
	S_Stack[s] = CallStack()
	return s, err
}

func (c *STConn) Close() error {
	if e := c.Ev.OnClose(c); e != nil {
		return e
	}
	delete(c.Dr.DBS, c.Db)
	if _, ok := C_Stack[c]; ok {
		delete(C_Stack, c)
	} else {
		fmt.Println("closing connectiong not found in static")
	}
	return c.Db.Close()
}

type STTx struct {
	Conn *STConn
	Tx   *sql.Tx
	Ev   DbEv //database evernt
}

func (tx *STTx) Commit() error {
	defer tx.Conn.TxDone()
	if e := tx.Ev.OnTxCommit(tx); e != nil {
		tx.Tx.Rollback()
		return e
	} else {
		return tx.Tx.Commit()
	}
}

func (tx *STTx) Rollback() error {
	defer tx.Conn.TxDone()
	if e := tx.Ev.OnTxRollback(tx); e != nil {
		tx.Tx.Rollback()
		return e
	} else {
		return tx.Tx.Rollback()
	}
}

type STStmt struct {
	Q    string
	Conn *STConn
	Stmt *sql.Stmt
	Ev   DbEv //database evernt
	Num  int
}

func (s *STStmt) NumInput() int {
	return s.Num
}

func (s *STStmt) Query(args []driver.Value) (driver.Rows, error) {
	if e := s.Ev.OnStmQuery(s, args); e != nil {
		return nil, e
	}
	targs := []interface{}{}
	for _, v := range args {
		targs = append(targs, v)
	}
	rows, e := s.Stmt.Query(targs...)
	rows_s := &STRows{
		Stmt: s,
		Rows: rows,
		Args: args,
		Ev:   s.Ev,
	}
	R_Stack[rows_s] = CallStack()
	return rows_s, e
}

func (s *STStmt) Exec(args []driver.Value) (driver.Result, error) {
	if e := s.Ev.OnStmExec(s, args); e != nil {
		return nil, e
	}
	targs := []interface{}{}
	for _, v := range args {
		targs = append(targs, v)
	}
	res, e := s.Stmt.Exec(targs...)
	res_s := &STResult{
		Stmt: s,
		Res:  res,
		Args: args,
		Ev:   s.Ev,
	}
	return res_s, e
}

func (s *STStmt) Close() error {
	if e := s.Ev.OnStmClose(s); e != nil {
		return e
	}
	if _, ok := S_Stack[s]; ok {
		delete(S_Stack, s)
	} else {
		fmt.Println("closing STMT not found in stack")
	}
	return s.Stmt.Close()
}

type STResult struct {
	Stmt *STStmt
	Res  sql.Result
	Args []driver.Value
	Ev   DbEv //database evernt
}

func (r *STResult) LastInsertId() (int64, error) {
	if e := r.Ev.OnResLIId(r); e != nil {
		return 0, e
	}
	return r.Res.LastInsertId()
}

func (r *STResult) RowsAffected() (int64, error) {
	if e := r.Ev.OnResARow(r); e != nil {
		return 0, e
	}
	return r.Res.RowsAffected()
}

type STRows struct {
	Stmt *STStmt
	Args []driver.Value
	Rows *sql.Rows
	Ev   DbEv //database evernt
}

func (rc *STRows) Columns() []string {
	cls, _ := rc.Rows.Columns()
	return cls
}

func (rc *STRows) Next(dest []driver.Value) error {
	if e := rc.Ev.OnRowNext(rc); e != nil {
		return e
	}
	if rc.Ev.IsEmpty(rc) {
		return io.EOF
	}
	if rc.Rows.Next() {
		l := len(rc.Columns())
		sary := make([]interface{}, l) //scan array.
		for i := 0; i < l; i++ {
			var a interface{}
			sary[i] = &a
		}
		e := rc.Rows.Scan(sary...)
		for i := 0; i < l; i++ {
			dest[i] = reflect.Indirect(reflect.ValueOf(sary[i])).Interface()
		}
		return e
	} else {
		return io.EOF
	}
}

func (rc *STRows) Close() error {
	if e := rc.Ev.OnRowClose(rc); e != nil {
		return e
	}
	if _, ok := R_Stack[rc]; ok {
		delete(R_Stack, rc)
	} else {
		fmt.Println("closing ROW not found in stack")
	}
	return rc.Rows.Close()
}

////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////// Log /////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////

var s_log bool = false

func ShowLog(v bool) {
	s_log = v
}
func log(f string, args ...interface{}) {
	if s_log {
		fmt.Println(fmt.Sprintf("DEM %v", fmt.Sprintf(f, args...)))
	}
}

////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////// Log /////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////

func OpenDem() *sql.DB {
	db, _ := sql.Open("DEM", G_Dsn)
	return db
}

//global databse name and connection
var G_Dn, G_Dsn string

//globack event callback instance.
var Evb *EcEv

//

//the type of TDbErr
type STErr uint32

//all error type
const (
	OPEN_ERR STErr = 1 << iota
	BEGIN_ERR
	CLOSE_ERR
	PREPARE_ERR
	TX_ROLLBACK_ERR
	TX_COMMIT_ERR
	STMT_CLOSE_ERR
	STMT_QUERY_ERR
	STMT_EXEC_ERR
	ROWS_CLOSE_ERR
	ROWS_NEXT_ERR
	EMPTY_DATA_ERR
	LAST_INSERT_ID_ERR
	ROWS_AFFECTED_ERR
)

func (t STErr) String() string {
	switch t {
	case OPEN_ERR:
		return "OPEN_ERR"
	case BEGIN_ERR:
		return "CONN_BEGIN_ERR"
	case CLOSE_ERR:
		return "CONN_CLOSE_ERR"
	case PREPARE_ERR:
		return "PREPARE_ERR"
	case TX_ROLLBACK_ERR:
		return "ROLLBACK_ERR"
	case TX_COMMIT_ERR:
		return "COMMIT_ERR"
	case STMT_CLOSE_ERR:
		return "STMT_CLOSE_ERR"
	case STMT_QUERY_ERR:
		return "STMT_QUERY_ERR"
	case STMT_EXEC_ERR:
		return "STMT_EXEC_ERR"
	case ROWS_CLOSE_ERR:
		return "ROWS_CLOSE_ERR"
	case ROWS_NEXT_ERR:
		return "ROWS_NEXT_ERR"
	case EMPTY_DATA_ERR:
		return "EMPTY_DATA_ERR"
	case LAST_INSERT_ID_ERR:
		return "LAST_INSERT_ID_ERR"
	case ROWS_AFFECTED_ERR:
		return "ROWS_AFFECTED_ERR"
	}
	return ""
}

//if error contain target error.
func (t STErr) Is(e STErr) bool {
	if (t & e) == e {
		return true
	} else {
		return false
	}
}

func (t STErr) IsErr(e STErr) error {
	if t.Is(e) {
		return errors.New(fmt.Sprintf("DEM %v", e.String()))
	} else {
		return nil
	}
}

type Query struct {
	Q    *regexp.Regexp
	Args *regexp.Regexp
}

func (q *Query) Match(query string, args []driver.Value) bool {
	return q.Q.MatchString(query) && q.Args.MatchString(fmt.Sprintf("%v", args))
}

//base event inteface implementation.
type EvBase struct {
	Errs  STErr
	QErr  []Query
	Dn    string
	IsErr func(STErr) error
}

func NewEvBase(dn string) *EvBase {
	eb := &EvBase{}
	eb.Dn = dn
	eb.IsErr = eb.Errs.IsErr
	return eb
}
func (e *EvBase) ResetErr() {
	e.SetErrs(0)
	e.ClsQErr()
}
func (e *EvBase) SetErrs(err STErr) {
	e.Errs = err
}
func (e *EvBase) AddErrs(err STErr) *EvBase {
	e.Errs = e.Errs | err
	return e
}
func (e *EvBase) ClsQErr() {
	e.QErr = []Query{}
}
func (e *EvBase) AddQErr(err Query) *EvBase {
	e.QErr = append(e.QErr, err)
	return e
}
func (e *EvBase) AddQErr2(qreg string, areg string) *EvBase {
	return e.AddQErr(Query{
		Q:    regexp.MustCompile(qreg),
		Args: regexp.MustCompile(areg),
	})
}
func (e *EvBase) AddQErr3(qreg string) *EvBase {
	return e.AddQErr2(qreg, ".*")
}
func (e *EvBase) Match(query string, args []driver.Value) bool {
	for _, q := range e.QErr {
		if q.Match(query, args) {
			return true
		}
	}
	return false
}
func (e *EvBase) OnOpen(dsn string) (*sql.DB, error) {
	err := e.IsErr(OPEN_ERR)
	if err != nil {
		return nil, err
	}
	dn := ""
	if len(e.Dn) < 1 {
		dn = G_Dn
	}
	if len(dn) < 1 {
		return nil, errors.New("dbname is not initial for event handler")
	}
	return sql.Open(dn, dsn)
}
func (e *EvBase) OnBegin(c *STConn) error {
	return e.IsErr(BEGIN_ERR)
}
func (e *EvBase) OnPrepare(c *STConn, query string) error {
	return e.IsErr(PREPARE_ERR)
}
func (e *EvBase) OnNumInput(c *STConn, query string, stm *sql.Stmt) int {
	return strings.Count(query, "?")
}
func (e *EvBase) OnClose(c *STConn) error {
	return e.IsErr(CLOSE_ERR)
}
func (e *EvBase) OnTxCommit(tx *STTx) error {
	return e.IsErr(TX_COMMIT_ERR)
}
func (e *EvBase) OnTxRollback(tx *STTx) error {
	return e.IsErr(TX_ROLLBACK_ERR)
}
func (e *EvBase) OnStmQuery(stm *STStmt, args []driver.Value) error {
	log("Query(%v) args(%v)", stm.Q, args)
	err := e.IsErr(STMT_QUERY_ERR)
	if err != nil {
		return err
	}
	if e.Match(stm.Q, args) {
		return errors.New("DEM query matched error")
	}
	return nil
}
func (e *EvBase) OnStmExec(stm *STStmt, args []driver.Value) error {
	log("Exec(%v) args(%v)", stm.Q, args)
	err := e.IsErr(STMT_EXEC_ERR)
	if err != nil {
		return err
	}
	if e.Match(stm.Q, args) {
		return errors.New("DEM query matched error")
	}
	return nil
}
func (e *EvBase) OnStmClose(stm *STStmt) error {
	return e.IsErr(STMT_CLOSE_ERR)
}
func (e *EvBase) OnResLIId(res *STResult) error {
	return e.IsErr(LAST_INSERT_ID_ERR)
}
func (e *EvBase) OnResARow(res *STResult) error {
	return e.IsErr(ROWS_AFFECTED_ERR)
}
func (e *EvBase) OnRowNext(row *STRows) error {
	return e.IsErr(ROWS_NEXT_ERR)
}
func (e *EvBase) IsEmpty(row *STRows) bool {
	return e.Errs.Is(EMPTY_DATA_ERR)
}
func (e *EvBase) OnRowClose(row *STRows) error {
	return e.IsErr(ROWS_CLOSE_ERR)
}

type EcEv struct {
	TC map[STErr]int
	MC map[STErr]int
	EvBase
}

func (e *EcEv) ResetErr() {
	e.EvBase.ResetErr()
	e.ClsEc()
}
func (e *EcEv) AddEC(err STErr, c int) *EcEv {
	e.TC[err] = c
	return e
}
func (e *EcEv) Match(err STErr) error {
	e.MC[err] = e.MC[err] + 1
	for k, v := range e.TC {
		if k.Is(err) && v == e.MC[err] {
			return errors.New(fmt.Sprintf("%v count(%v) error", err.String(), v))
		}
	}
	return nil
}
func (e *EcEv) CheckErr(err STErr) error {
	err_m := e.Errs.IsErr(err)
	if err_m == nil {
		return e.Match(err)
	} else {
		return err_m
	}
}
func (e *EcEv) ClsEc() {
	e.TC = map[STErr]int{}
	e.MC = map[STErr]int{}
}

func NewEcEv(dn string) *EcEv {
	evec := &EcEv{}
	evec.TC = map[STErr]int{}
	evec.MC = map[STErr]int{}
	evec.Dn = dn
	evec.IsErr = evec.CheckErr
	return evec
}

func init() {
	Evb = NewEcEv("")
	Register("DEM", Evb)
}
