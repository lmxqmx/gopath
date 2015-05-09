package sr

import (
	"fmt"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

type srh struct {
	c int64
}

func (d *srh) Path(hs *routing.HTTPSession, sr *SR) (string, error) {
	if d.c < 1 {
		d.c++
		return "", util.Err("normal err")
	} else {
		return fmt.Sprintf("%v/%v-%v", sr.R, util.Now(), atomic.AddInt64(&d.c, 1)), nil
	}
}

func (d *srh) OnSrF(hs *routing.HTTPSession, aid, ver, dev, sp, sf string) error {
	return util.Err("normal err")
}
func (d *srh) OnSrL(hs *routing.HTTPSession, aid, ver, prev, dev string, from, all int64) (interface{}, error) {
	return nil, util.Err("normal error")
}
func (d *srh) OnSrPkg(hs *routing.HTTPSession, aid, dev string) (interface{}, error) {
	return map[string]interface{}{}, nil
}

type srh_q_h struct {
	b   bool
	le  int
	pkg bool
}

func (sr *srh_q_h) Args(s *SRH_Q, hs *routing.HTTPSession, aid, ver, dev, sp, sf string) (util.Map, error) {
	return hs.AllRVal(), nil
}
func (sr *srh_q_h) Proc(s *SRH_Q, i *SRH_Q_I) error {
	if sr.b {
		return util.Err("normal error")
	}
	for _, ev := range i.Evs {
		fmt.Println(ev)
	}
	return nil
}
func (sr *srh_q_h) ListSr(s *SRH_Q, hs *routing.HTTPSession, aid, ver, prev, dev string, from, all int64) (interface{}, error) {
	switch sr.le {
	case 1:
		return nil, util.Err("normal error")
	case 2:
		return nil, util.NOT_FOUND
	default:
		return []interface{}{}, nil
	}
}
func (sr *srh_q_h) ListPkg(s *SRH_Q, hs *routing.HTTPSession, aid, dev string) (interface{}, error) {
	if sr.pkg {
		return nil, util.Err("normal error")
	} else {
		sr.pkg = true
		return map[string]interface{}{}, nil
	}
}
func TestAddSr(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	sr := NewSR("/tmp")
	ts := httptest.NewServer2(sr)
	ts.G("")
	ts.G("?aid=sss&ver=1.0.0")
	for i := 0; i < 10; i++ {
		ts.PostF2("", "sr_f", "sr_f.zip", map[string]string{
			"exec": "A",
			"aid":  "org..",
			"ver":  "0.0.1",
		})
	}
	ts.PostF2("", "sr_f", "sr.go", map[string]string{
		"exec": "A",
		"aid":  "org..",
		"ver":  "0.0.1",
	})
	ts.PostF2("", "sr_f", "sr.go", map[string]string{
		"exec": "A",
	})
	sr.R = "/sdfs"
	ts.PostF2("", "sr_f", "sr.go", map[string]string{
		"exec": "A",
		"aid":  "org..",
		"ver":  "0.0.1",
	})
	sr.R = "/tmp"
	sr.H = &srh{}
	ts.PostF2("", "sr_f", "sr_f.zip", map[string]string{
		"exec": "A",
		"aid":  "org..",
		"ver":  "0.0.1",
	})
	ts.PostF2("", "sr_f", "sr_f.zip", map[string]string{
		"exec": "A",
		"aid":  "org..",
		"ver":  "0.0.1",
	})

	//
	sqh := &srh_q_h{}
	sr2, srh_q := NewSR3("/tmp", sqh)
	ts2 := httptest.NewServer2(sr2)
	ts2.PostF2("", "sr_f", "sr_f.zip", map[string]string{
		"exec": "A",
		"aid":  "org..",
		"ver":  "0.0.1",
	})
	srh_q.Run(5)
	for i := 0; i < 10; i++ {
		ts2.PostF2("", "sr_f", "sr_f.zip", map[string]string{
			"exec": "A",
			"aid":  "org..",
			"ver":  "0.0.1",
		})
	}
	fmt.Println(ts2.PostF2("", "sr_f", "sr.go", map[string]string{
		"exec": "A",
		"aid":  "org..",
		"ver":  "0.0.1",
	}))
	fmt.Println("---->")
	ts2.PostF2("", "sr_f", "sr.zip", map[string]string{
		"exec": "A",
		"aid":  "org..",
		"ver":  "0.0.1",
	})
	util.FWrite2("er.dat", []byte{0, 0, 'a', 'b', 'c'})
	util.Zip("er.zip", ".", "./er.dat")
	ts2.PostF2("", "sr_f", "er.zip", map[string]string{
		"exec": "A",
		"aid":  "org..",
		"ver":  "0.0.1",
	})
	fmt.Println(ts2.G2("?exec=L&aid=org&ver=0.0.1"))
	sqh.le = 2
	fmt.Println(ts2.G2("?exec=L&aid=org&ver=0.0.1"))
	sqh.le = 1
	fmt.Println(ts2.G2("?exec=L&aid=org&ver=0.0.1"))
	fmt.Println(ts2.G2("?exec=L&aid=org&ver=0.0.1&prev="))
	fmt.Println(ts2.G2("?exec=L&pn=sss"))
	fmt.Println(ts2.G2("?exec=L&aid=org&ver=0.0.1&all=x"))
	fmt.Println(ts2.G2("?exec=L&aid=org&ver=0.0.1&from=999999999999"))
	fmt.Println(ts2.G2("?exec=P"))
	fmt.Println(ts2.G2("?exec=P"))
	fmt.Println(ts2.G2("?exec=X"))
	time.Sleep(500 * time.Millisecond)
	sqh.b = true
	ts2.PostF2("", "sr_f", "sr_f.zip", map[string]string{
		"exec": "A",
		"aid":  "org..",
		"ver":  "0.0.1",
	})
	time.Sleep(500 * time.Millisecond)
	srh_q.Stop()
	time.Sleep(time.Second)
	//
	kk := &SRH_N{}
	kk.OnSrL(nil, "aid", "ver", "", "", 0, 0)
	kk.OnSrPkg(nil, "aid", "dev")
}

// func TestListSr(t *testing.T) {
// 	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
// 	sqh := &srh_q_h{}
// 	sr, srh_q := NewSR3("/tmp", sqh)
// 	srh_q.Run(5)
// 	ts := httptest.NewServer(sr.ListSr)
// 	fmt.Println(ts.G2("?aid=org&ver=0.0.1"))
// 	sqh.le = true
// 	fmt.Println(ts.G2("?aid=org&ver=0.0.1"))
// 	fmt.Println(ts.G2(""))
// 	kk := &SRH_N{}
// 	kk.OnSrL(nil, "aid", "ver", 0, 0, 0)
// }
