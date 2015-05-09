package sr

import (
	"bufio"
	"fmt"
	"github.com/Centny/gwf/iow"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/sr/pb"
	"github.com/Centny/gwf/util"
	"github.com/golang/protobuf/proto"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

type SRH interface {
	Path(hs *routing.HTTPSession, sr *SR) (string, error)
	OnSrF(hs *routing.HTTPSession, aid, ver, dev, sp, sf string) error
	OnSrL(hs *routing.HTTPSession, aid, ver, prev, dev string, from, all int64) (interface{}, error)
	OnSrPkg(hs *routing.HTTPSession, aid, dev string) (interface{}, error)
}
type SR struct {
	H SRH
	R string //root store path.
}

func NewSR(r string) *SR {
	return &SR{
		R: r,
		H: &SRH_N{
			c: 0,
		},
	}
}
func NewSR2(r string, h SRH) *SR {
	return &SR{
		R: r,
		H: h,
	}
}
func NewSR3(r string, h SRH_Q_H) (*SR, *SRH_Q) {
	sq := NewSRH_Q(r, h)
	return NewSR2(r, sq), sq
}
func (s *SR) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var action string = "L"
	err := hs.ValidCheckVal(`
		exec,O|S,O:A~L~P;
		`, &action)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	switch action {
	case "A":
		return s.AddSr(hs)
	case "P":
		return s.ListPkg(hs)
	default:
		return s.ListSr(hs)
	}
}
func (s *SR) AddSr(hs *routing.HTTPSession) routing.HResult {
	var dev, aid, ver string
	err := hs.ValidCheckVal(`
		dev,O|S,L:0;
		aid,R|S,L:0;
		ver,R|S,L:0;
		`, &dev, &aid, &ver)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	sp, err := s.H.Path(hs, s)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	sf := fmt.Sprintf("%v/%v/sr_f.zip", s.R, sp)
	_, err = hs.RecF("sr_f", sf)
	if err != nil {
		return hs.MsgResErr2(1, "srv-err", err)
	}
	err = s.H.OnSrF(hs, aid, ver, dev, sp, sf)
	if err == nil {
		return hs.MsgRes("OK")
	} else {
		return hs.MsgResErr2(1, "srv-err", err)
	}
}
func (s *SR) ListSr(hs *routing.HTTPSession) routing.HResult {
	var prev string
	var from, all int64 = 0, 0
	var dev, aid, ver string
	err := hs.ValidCheckVal(`
		prev,O|S,L:0;
		from,O|I,R:0;
		all,O|I,O:0~1;
		dev,O|S,L:0;
		aid,R|S,L:0;
		ver,R|S,L:0;
		`, &prev, &from, &all, &dev, &aid, &ver)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	data, err := s.H.OnSrL(hs, aid, ver, prev, dev, from, all)
	if err == nil {
		return hs.MsgRes(data)
	} else if err == util.NOT_FOUND {
		return hs.MsgRes2(404, []interface{}{})
	} else {
		return hs.MsgResErr2(1, "srv-err", err)
	}
}

func (s *SR) ListPkg(hs *routing.HTTPSession) routing.HResult {
	var dev, aid string
	hs.ValidCheckVal(`
		dev,O|S,L:0;
		aid,O|S,L:0;
		`, &dev, &aid)
	data, err := s.H.OnSrPkg(hs, aid, dev)
	if err == nil {
		return hs.MsgRes(data)
	} else {
		return hs.MsgResErr2(1, "srv-err", err)
	}
}

type SRH_N struct {
	c int64
}

func (s *SRH_N) Path(hs *routing.HTTPSession, sr *SR) (string, error) {
	return fmt.Sprintf("%v-%v", util.Now(), atomic.AddInt64(&s.c, 1)), nil
}

func (s *SRH_N) OnSrF(hs *routing.HTTPSession, aid, ver, dev, sp, sf string) error {
	return nil
}
func (s *SRH_N) OnSrL(hs *routing.HTTPSession, aid, ver, prev, dev string, from, all int64) (interface{}, error) {
	return []interface{}{}, nil
}
func (s *SRH_N) OnSrPkg(hs *routing.HTTPSession, aid, dev string) (interface{}, error) {
	return map[string]interface{}{}, nil
}

type SRH_Q_I struct {
	Id   interface{} `bson:"_id" json:"id"`
	Sp   string      `json:"sp"`
	Aid  string      `json:"aid"`
	Ver  string      `json:"ver"`
	Dev  string      `json:"dev"`
	Kvs  util.Map    `json:"-" bson:"-"`
	Evs  []*pb.Evn   `json:"evs"`
	Time int64       `json:"time"`
}
type SRH_Q_H interface {
	Args(s *SRH_Q, hs *routing.HTTPSession, aid, ver, dev, sp, sf string) (util.Map, error)
	Proc(s *SRH_Q, i *SRH_Q_I) error
	ListSr(s *SRH_Q, hs *routing.HTTPSession, aid, ver, prev, dev string, from, all int64) (interface{}, error)
	ListPkg(s *SRH_Q, hs *routing.HTTPSession, aid, dev string) (interface{}, error)
}
type SRH_Q struct {
	SRH_N
	R       string
	H       SRH_Q_H
	Q       chan *SRH_Q_I
	Running bool
}

func NewSRH_Q(r string, h SRH_Q_H) *SRH_Q {
	return &SRH_Q{
		R: r,
		H: h,
		Q: make(chan *SRH_Q_I, 3000),
	}
}
func (s *SRH_Q) OnSrF(hs *routing.HTTPSession, aid, ver, dev, sp, sf string) error {
	if !s.Running {
		log.W("SRH_Q OnSrF err:Proc is not running")
		return util.Err("SRH_Q not running")
	}
	kvs, err := s.H.Args(s, hs, aid, ver, dev, sp, sf)
	if err == nil {
		s.Q <- &SRH_Q_I{
			Sp:  sp,
			Aid: aid,
			Ver: ver,
			Dev: dev,
			Kvs: kvs,
		}
	}
	return err
}
func (s *SRH_Q) OnSrL(hs *routing.HTTPSession, aid, ver, prev, dev string, from, all int64) (interface{}, error) {
	return s.H.ListSr(s, hs, aid, ver, prev, dev, from, all)
}
func (s *SRH_Q) OnSrPkg(hs *routing.HTTPSession, aid, dev string) (interface{}, error) {
	return s.H.ListPkg(s, hs, aid, dev)
}
func (s *SRH_Q) Proc() {
	tick := time.Tick(500 * time.Millisecond)
	for s.Running {
		select {
		case i := <-s.Q:
			s.doproc(i)
		case <-tick:
		}
	}
	log.D("SRH_Q Proc done...")
}
func (s *SRH_Q) doproc(i *SRH_Q_I) {
	sr_p := filepath.Join(s.R, i.Sp)
	sr_f := filepath.Join(s.R, i.Sp, "sr_f.zip")
	err := util.Unzip(sr_f, sr_p)
	if err != nil {
		log.E("unzip %v err:%v", sr_f, err.Error())
		return
	}
	sr_er := filepath.Join(sr_p, "er.dat")
	er_f, err := os.Open(sr_er)
	if err != nil {
		log.E("open er.data file %v err:%v", sr_er, err.Error())
		return
	}
	err = iow.ReadLdata(bufio.NewReader(er_f), func(bys []byte) error {
		var evn pb.Evn
		err = proto.Unmarshal(bys, &evn)
		if err == nil {
			i.Evs = append(i.Evs, &evn)
		}
		return err
	})
	er_f.Close()
	if err != nil && err != io.EOF {
		log.E("Unmarshal er.data file %v err:%v", sr_er, err.Error())
		return
	}
	i.Time = util.Now()
	err = s.H.Proc(s, i)
	if err == nil {
		log.D("Proc SRH_Q_I(%v) OK", i.Id)
	} else {
		log.E("Proc SRH_Q_I %v err:%v", i, err.Error())
	}
}
func (s *SRH_Q) Run(c int) {
	s.Running = true
	for i := 0; i < c; i++ {
		go s.Proc()
	}
	log.I("SRH_Q Run %v Proc", c)
}
func (s *SRH_Q) Stop() {
	s.Running = false
	log.I("SRH_Q Stopping Proc")
}
