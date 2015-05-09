//Package netw provide the base transfer protocol for TCP
//
//it contain the client and server base struct that can be extended by event handler.
//
//Protocol:mod->len->data
package netw

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/pool"
	"github.com/Centny/gwf/util"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var con_idc uint64 = 0
var pool_idc uint64 = 0

//whether show debug log or not,default is false.
var ShowLog bool = false
var ShowLog_C bool = false

func log_d(f string, args ...interface{}) {
	if ShowLog {
		log.D_(1, f, args...)
	}
}

//the protocol modes
const H_MOD = "^~^"

//the default connection timeout for not data receive
const CON_TIMEOUT int64 = 5000

//the func to create on connection for ConPool
type NewConF func(cp ConPool, p *pool.BytePool, con net.Conn) Con

//the func to handler Cmd.Err call.
type CmdErrF func(c Cmd, d int, code byte, f string, args ...interface{})

//func for covert struct to []byte
type V2Byte func(v interface{}) ([]byte, error)

//func for covert []byte to struct or map.
type Byte2V func(bys []byte, v interface{}) (interface{}, error)

//the base connection event handler
type ConHandler interface {
	//calling when the connection have been connected.
	OnConn(c Con) bool
	//calling when the connection have been closed.
	OnClose(c Con)
}

//the connect data event handler.
type CmdHandler interface {
	//calling when one entire command have been received.
	OnCmd(c Cmd) int
}

//all ConPool handler contain CmdHandler and ConHandler.
type CCHandler interface {
	ConHandler
	CmdHandler
}

/*

*/

//the connection interface.
//it will be created when client connect to server or server received one connection.
type Con interface {
	//the base connection
	net.Conn
	//the ConPool
	CP() ConPool
	//the memory pool
	P() *pool.BytePool
	// R() *bufio.Reader  //the buffer reader
	// W() *bufio.Writer  //the buffer writer.
	//the connection id.
	Id() string
	//set the connection id.
	SetId(id string)
	//connection seesion.
	Kvs() util.Map
	//the last update time for data transfer
	Last() int64
	//set connection wait status, if true,the connection will not timeout
	SetWait(t bool)
	//if connection in waiting status.
	Waiting() bool
	//read byte data and wait until have receive p length data.
	ReadW(p []byte) error
	//write multi []byte to conection.
	//it will be joined to MOD|lenght|[]byte|[]byte|[]byte....
	Writeb(bys ...[]byte) (int, error)
	//write one struct val to connection.
	//it will call connection V2B func to convert the value to []byte.
	Writev(val interface{}) (int, error)
	//exec on remote command by args,
	//the return value will be converted to dest,and return dest
	Exec(args interface{}, dest interface{}) (interface{}, error)
	//flush the buffer.
	Flush() error
	//the value to []byte convert function
	V2B() V2Byte
	//the []byte to value convert function
	B2V() Byte2V
}

//the base implement to Con
type Con_ struct {
	net.Conn                //the base connection
	CP_      ConPool        //the ConPool
	P_       *pool.BytePool //the memory pool
	R_       *bufio.Reader  //the buffer reader
	W_       *bufio.Writer  //the buffer writer.
	Kvs_     util.Map       //the session.
	Last_    int64          //the last update time for data transfer
	Waiting_ int32          //whether in waiting status.
	V2B_     V2Byte         //the V2Byte func
	B2V_     Byte2V         //the Byte2V func
	ID_      string         //the connection id
	c_l      sync.RWMutex   //connection lock.
	buf      []byte         //the buf to store the data len which will be writed to connection.
	ShowLog  bool
	// r_l      sync.RWMutex
}

//new Con by ConPool/BytePool and normal connection.
func NewCon(cp ConPool, p *pool.BytePool, con net.Conn) Con {
	return NewCon_(cp, p, con)
}

//Con_ creator.
func NewCon_(cp ConPool, p *pool.BytePool, con net.Conn) *Con_ {
	return &Con_{
		CP_:      cp,
		P_:       p,
		Conn:     con,
		R_:       bufio.NewReader(con),
		W_:       bufio.NewWriter(con),
		Kvs_:     util.Map{},
		Waiting_: 0,
		buf:      make([]byte, 2),
		V2B_: func(v interface{}) ([]byte, error) {
			return nil, util.Err("V2B not implemeted")
		},
		B2V_: func(bys []byte, v interface{}) (interface{}, error) {
			return nil, util.Err("B2V not implemeted")
		},
		ID_:     fmt.Sprintf("C%v", atomic.AddUint64(&con_idc, 1)),
		ShowLog: ShowLog_C,
	}
}
func (c *Con_) log_d(f string, args ...interface{}) {
	if c.ShowLog {
		log.D(f, args...)
	}
}
func (c *Con_) CP() ConPool {
	return c.CP_
}
func (c *Con_) P() *pool.BytePool {
	return c.P_
}

// func (c *Con_) R() *bufio.Reader {
// 	return c.R_
// }
// func (c *Con_) W() *bufio.Writer {
// 	return c.W_
// }
//

func (c *Con_) Kvs() util.Map {
	return c.Kvs_
}
func (c *Con_) Last() int64 {
	return c.Last_
}

//set the connection waiting status.
//if true,the connection will keep forever.
//if false,the connection will be closed after timeout when not data receive.
func (c *Con_) SetWait(t bool) {
	if t {
		atomic.StoreInt32(&c.Waiting_, 1)
	} else {
		atomic.StoreInt32(&c.Waiting_, 0)
	}
}
func (c *Con_) Waiting() bool {
	return c.Waiting_ > 0
}

//read the number of the data in p
func (c *Con_) ReadW(p []byte) error {
	// c.r_l.Lock()
	// defer c.r_l.Unlock()
	return util.ReadW(c.R_, p, &c.Last_)
}

//rewrite to forbiden call.
func (c *Con_) Write(b []byte) (n int, err error) {
	panic("do not call Write direct,using Writeb/Writev instead")
}

//sending data.
//Data:mod|len|bys...
func (c *Con_) Writeb(bys ...[]byte) (int, error) {
	c.c_l.Lock()
	defer c.c_l.Unlock()
	total, _ := Writeb(c.W_, bys...)
	c.log_d("write data(%v) to %v", total, c.RemoteAddr().String())
	return total, c.Flush()
}
func (c *Con_) Writev(val interface{}) (int, error) {
	return Writev(c, val)
}
func (c *Con_) Exec(args interface{}, dest interface{}) (interface{}, error) {
	return nil, util.Err("connection not implement Exec")
}
func (c *Con_) Flush() error {
	return c.W_.Flush()
}
func (c *Con_) V2B() V2Byte {
	return c.V2B_
}
func (c *Con_) B2V() Byte2V {
	return c.B2V_
}
func (c *Con_) Id() string {
	return c.ID_
}
func (c *Con_) SetId(id string) {
	c.ID_ = id
}

//the Cmd interface for exec data.
type Cmd interface {
	//get the connect.
	Con
	//get the base connection.
	BaseCon() Con
	//get the command data.
	Data() []byte
	//done the command, the data []byte will free.
	Done()
	//convert the data to dest value.
	//it will call the connection B2V func
	V(dest interface{}) (interface{}, error)
	//the error log stack depth.
	SetErrd(d int)
	//common error executor
	Err(code byte, f string, args ...interface{})
}

//the base implement to Cmd interface.
type Cmd_ struct {
	Con          //base connection.
	Data_ []byte //received data
	data_ []byte //really address to free.
	d     int    //the error log stack depth.
}

func (c *Cmd_) BaseCon() Con {
	return c.Con
}
func (c *Cmd_) Data() []byte {
	return c.Data_
}

//free the memory(Data []byte)
func (c *Cmd_) Done() {
	c.P().Free(c.data_)
}
func (c *Cmd_) V(dest interface{}) (interface{}, error) {
	return V(c, dest)
}
func (c *Cmd_) SetErrd(d int) {
	c.d = d
}
func (c *Cmd_) Err(code byte, f string, args ...interface{}) {
	c.CP().Err()(c, c.d, code, f, args...)
}

type ConPool interface {
	LoopTimeout()
	Close()
	RunC(c Con)
	Err() CmdErrF
	Find(id string) Con
	Id() string
	SetId(id string)
}

//the connection pool
type LConPool struct {
	T  int64          //the timeout of not data received
	P  *pool.BytePool //the memory pool
	Wg sync.WaitGroup //wait group.
	H  CCHandler      //command handler
	// Wc     chan int       //the wait chan.
	NewCon NewConF
	t_r    bool
	cons   map[string]Con
	cons_l sync.RWMutex
	Err_   CmdErrF
	Id_    string
}

//new connection pool.
func NewLConPool(p *pool.BytePool, h CCHandler, n string) *LConPool {
	return &LConPool{
		T:      CON_TIMEOUT,
		P:      p,
		H:      h,
		cons:   map[string]Con{},
		NewCon: NewCon,
		Err_: func(c Cmd, d int, code byte, f string, args ...interface{}) {
			log.D_(d, f, args...)
		},
		Id_: fmt.Sprintf("%v%v", n, atomic.AddUint64(&pool_idc, 1)),
	}
}

//looping the connection timeout.
func (l *LConPool) LoopTimeout() {
	l.t_r = true
	for l.t_r {
		cons := []Con{}
		tn := util.Now()
		for _, c := range l.cons {
			if c.Waiting() {
				continue
			}
			if (tn - c.Last()) > l.T {
				cons = append(cons, c)
			}
		}
		if (len(cons)) > 0 {
			log.D("closing %v connection for timeout", len(cons))
		}
		for _, con := range cons {
			con.Close()
		}
		time.Sleep(time.Duration(l.T) * time.Millisecond)
	}
	// l.Wc <- 0
}

//close all connection
func (l *LConPool) Close() {
	l.t_r = false
	l.cons_l.Lock()
	for _, c := range l.cons {
		c.Close()
	}
	l.cons = map[string]Con{}
	l.cons_l.Unlock()
}
func (l *LConPool) add_c(c Con) {
	l.cons_l.Lock()
	defer l.cons_l.Unlock()
	if _, ok := l.cons[c.Id()]; ok {
		panic(fmt.Sprintf("conection by id(%v) already added", c.Id()))
	}
	l.Wg.Add(1)
	l.cons[c.Id()] = c
	log_d("add connect(%v) to pool(%v)", c.Id(), l.Id())
}
func (l *LConPool) del_c(c Con) {
	l.cons_l.Lock()
	defer l.cons_l.Unlock()
	l.Wg.Done()
	delete(l.cons, c.Id())
	log_d("del connect(%v) from pool(%v)", c.Id(), l.Id())
}
func (l *LConPool) Find(id string) Con {
	if c, ok := l.cons[id]; ok {
		return c
	} else {
		return nil
	}
}

//run one connection by async.
func (l *LConPool) RunC(con Con) {
	// go func(lll *LConPool, conn net.Conn) {
	l.add_c(con) //will remove in RunC_ defer
	go l.runc_(con)
	// }(l, con)
}

//run on connection by sync.
func (l *LConPool) RunC_(con Con) {
	l.add_c(con)
	l.runc_(con)
}
func (l *LConPool) runc_(con Con) {
	defer func() {
		log_d("closing connection(%v,%v) in pool(%v)", con.RemoteAddr().String(), con.Id(), l.Id())
		l.H.OnClose(con)
		con.Close()
		l.del_c(con)
		if err := recover(); err != nil {
			buf := make([]byte, 102400)
			blen := runtime.Stack(buf, false)
			log.E("RunC_ close err(%v),stack:\n%v", err, string(buf[0:blen]))
		}
	}()
	log_d("running connection(%v,%v) in pool(%v)", con.RemoteAddr().String(), con.Id(), l.Id())
	//
	buf := make([]byte, 5)
	mod := []byte(H_MOD)
	mod_l := len(mod)
	//
	for {
		err := con.ReadW(buf)
		if err != nil {
			log_d("read head mod from(%v) error:%v", con.RemoteAddr().String(), err.Error())
			break
		}
		if !bytes.HasPrefix(buf, mod) {
			log.W("reading invalid mod(%v) from(%v)", string(buf), con.RemoteAddr().String())
			continue
		}
		dlen := binary.BigEndian.Uint16(buf[mod_l:])
		if dlen < 2 {
			log.W("reading invalid data len for mod(%v) from(%v)", string(buf), con.RemoteAddr().String())
			continue
		}
		dbuf := l.P.Alloc(int(dlen))
		err = con.ReadW(dbuf)
		if err != nil {
			log_d("read data from(%v) error:%v", con.RemoteAddr().String(), err.Error())
			break
		}
		// if len(dbuf) < 3 {
		// 	continue
		// }
		l.H.OnCmd(&Cmd_{
			Con:   con,
			Data_: dbuf,
			data_: dbuf,
			d:     2,
		})
	}
}
func (l *LConPool) Err() CmdErrF {
	return l.Err_
}

func (l *LConPool) Id() string {
	return l.Id_
}
func (l *LConPool) SetId(id string) {
	l.Id_ = id
}
func (l *LConPool) Cons() map[string]Con {
	return l.cons
}

// func (l *LConPool) Write(bys []byte) int {
// 	for _, con := range l.cons {
// 		con.Write(bys)
// 	}
// 	return len(l.cons)
// }
func (l *LConPool) Writeb(bys ...[]byte) int {
	for _, con := range l.cons {
		con.Writeb(bys...)
	}
	return len(l.cons)
}
func (l *LConPool) Writev(val interface{}) int {
	for _, con := range l.cons {
		con.Writev(val)
	}
	return len(l.cons)
}
