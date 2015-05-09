package im

// import (
// 	"fmt"
// 	"github.com/Centny/gwf/netw"
// 	"github.com/Centny/gwf/netw/impl"
// 	"github.com/Centny/gwf/util"
// 	"math/rand"
// 	"strings"
// 	"sync"
// 	"sync/atomic"
// 	"time"
// )

// type RCDbH struct {
// 	RC *impl.RC_Con
// }

// func NewRCDbH(rc *impl.RC_Con) *RCDbH {
// 	return &RCDbH{
// 		RC: rc,
// 	}
// }
// func (r *RCDbH) OnConn(c netw.Con) bool {
// 	r.RC.Execm(m, args, dest)
// 	return true
// }

// //calling when the connection have been closed.
// func (r *RCDbH) OnClose(c netw.Con) {
// }
// func (r *RCDbH) AddCon(c *Con) error {
// 	if c == nil {
// 		panic("Con is nil")
// 	}
// 	m.con_l.Lock()
// 	defer m.con_l.Unlock()
// 	m.Cons[fmt.Sprintf("%v%v%v%v", c.Sid, c.Cid, c.R, c.T)] = c
// 	return nil
// }
// func (r *RCDbH) DelCon(sid, cid, r string, t byte) error {
// 	m.con_l.Lock()
// 	defer m.con_l.Unlock()
// 	delete(m.Cons, fmt.Sprintf("%v%v%v%v", sid, cid, r, t))
// 	return nil
// }

// //list all connection by target R
// func (r *RCDbH) ListCon(rs []string) ([]Con, error) {
// 	if m == nil {
// 		panic(nil)
// 	}
// 	rsm := map[string]byte{}
// 	for _, r := range rs {
// 		rsm[r] = 1
// 	}
// 	ccs := []Con{}
// 	for _, cc := range m.Cons {
// 		if _, ok := rsm[cc.R]; ok {
// 			ccs = append(ccs, *cc)
// 		}
// 	}
// 	return ccs, nil
// }

// //
// //
// //list all user R by group R
// func (r *RCDbH) ListUsrR(gr []string) ([]string, error) {
// 	trs := []string{}
// 	for _, g := range gr {
// 		if rs, ok := m.Grp[g]; ok {
// 			trs = append(trs, rs...)
// 		}
// 	}
// 	return trs, nil
// }
// func (r *RCDbH) ListR() ([]string, error) {
// 	var usrs []string = []string{}
// 	for r, _ := range m.Usr {
// 		usrs = append(usrs, r)
// 	}
// 	return usrs, nil
// }

// //sift the R to group R and user R.
// func (r *RCDbH) Sift(rs []string) ([]string, []string, error) {
// 	ur, gr := []string{}, []string{}
// 	for _, r := range rs {
// 		if strings.HasPrefix(r, "G-") {
// 			gr = append(gr, r)
// 		} else {
// 			ur = append(ur, r)
// 		}
// 	}
// 	return gr, ur, nil
// }

// //
// //
// func (r *RCDbH) AddSrv(srv *Srv) error {
// 	m.srv_l.Lock()
// 	defer m.srv_l.Unlock()
// 	// srv.Token = "abc"
// 	// fmt.Println(m, srv)
// 	m.Srvs[srv.Sid] = srv
// 	return nil
// }
// func (r *RCDbH) DelSrv(sid string) error {
// 	m.srv_l.Lock()
// 	defer m.srv_l.Unlock()
// 	delete(m.Srvs, sid)
// 	return nil
// }

// //find the server by token
// func (r *RCDbH) FindSrv(token string) (*Srv, error) {
// 	for _, srv := range m.Srvs {
// 		if srv.Token == token {
// 			return srv, nil
// 		}
// 	}
// 	return nil, util.Err("server not found by token(%v)", token)
// }

// //list all online server,exclue special server id.
// func (r *RCDbH) ListSrv(sid string) ([]Srv, error) {
// 	srvs := []Srv{}
// 	// fmt.Println(m, m.Srvs)
// 	for _, srv := range m.Srvs {
// 		if len(sid) > 0 && srv.Sid == sid {
// 			continue
// 		}
// 		srvs = append(srvs, *srv)
// 	}
// 	return srvs, nil
// }

// //
// //
// //user login,return user R.
// func (r *RCDbH) OnUsrLogin(r netw.Cmd, args *util.Map) (string, error) {
// 	m.u_lck.Lock()
// 	defer m.u_lck.Unlock()
// 	if args.Exist("token") {
// 		ur := fmt.Sprintf("U-%v", atomic.AddUint64(&m.u_cc, 1))
// 		m.Usr[ur] = 1
// 		log_d("user login by R(%v)", ur)
// 		return ur, nil
// 	} else {
// 		log_d("user login fail for token not found")
// 		return "", util.Err("login fail:token not found")
// 	}
// }
// func (r *RCDbH) OnUsrLogout(r string, args *util.Map) error {
// 	m.u_lck.Lock()
// 	defer m.u_lck.Unlock()
// 	if _, ok := m.Usr[r]; ok {
// 		delete(m.Usr, r)
// 		log_d("user logout by R(%v)", r)
// 		return nil
// 	} else {
// 		log_d("user logout fail:R not found")
// 		return util.Err("login fail:R not found")
// 	}
// }

// //
// //
// //update the message R status
// func (r *RCDbH) Update(ms *Msg, rs map[string]string) error {
// 	m.ms_l.Lock()
// 	defer m.ms_l.Unlock()
// 	if tm, ok := m.Ms[ms.Id]; ok {
// 		for r, s := range rs {
// 			tm.Ms[r] = s
// 		}
// 		m.Ms[tm.Id] = tm
// 		return nil
// 	} else {
// 		return util.Err("message not found by id(%v)", ms.Id)
// 	}
// }

// //store mesage
// func (r *RCDbH) Store(ms *Msg) error {
// 	m.ms_l.Lock()
// 	defer m.ms_l.Unlock()
// 	ms.Id = fmt.Sprintf("M-%v", atomic.AddUint64(&m.m_cc, 1))
// 	m.Ms[ms.Id] = ms
// 	return nil
// }

// func (r *RCDbH) RandGrp() (string, int) {
// 	if len(m.Grp) < 1 {
// 		return "", 0
// 	}
// 	gs := []string{}
// 	for gr, _ := range m.Grp {
// 		gs = append(gs, gr)
// 	}
// 	g := gs[rand.Intn(len(gs))]
// 	return g, len(m.Grp[g])
// }
// func (r *RCDbH) RandUsr() []string {
// 	ulen := len(m.Usr)
// 	if ulen < 1 {
// 		return []string{}
// 	}
// 	usrs, _ := m.ListR()
// 	um := map[string]byte{}
// 	tlen := rand.Intn(ulen)%16 + 1
// 	for i := 0; i <= tlen; i++ {
// 		um[usrs[rand.Intn(ulen)]] = 1
// 	}
// 	tur := []string{}
// 	for u, _ := range um {
// 		tur = append(tur, u)
// 	}
// 	return tur
// }
// func (r *RCDbH) GrpBuilder() {
// 	for {
// 		time.Sleep(time.Second)
// 		if len(m.Usr) < 1 {
// 			continue
// 		}
// 		usrs, _ := m.ListR()
// 		g := fmt.Sprintf("G-%v", atomic.AddUint64(&m.g_cc, 1))
// 		us := []string{}
// 		tlen := rand.Intn(len(m.Usr)) + 1
// 		mu := map[string]bool{}
// 		for i := 0; i < tlen; i++ {
// 			mu[usrs[rand.Intn(len(m.Usr))]] = true
// 		}
// 		for u, _ := range mu {
// 			us = append(us, u)
// 		}
// 		m.Grp[g] = us
// 	}
// }
// func (r *RCDbH) Show() (uint64, uint64, uint64, uint64, uint64) {
// 	mlen := uint64(len(m.Ms))
// 	var rlen uint64 = 0
// 	var plen uint64 = 0
// 	var elen uint64 = 0
// 	var dlen uint64 = 0
// 	for _, m := range m.Ms {
// 		rlen += uint64(len(m.Ms))
// 		for _, s := range m.Ms {
// 			if strings.HasPrefix(s, "E-") {
// 				elen++
// 			} else if s == MS_PENDING {
// 				plen++
// 			} else {
// 				dlen++
// 			}
// 		}
// 	}
// 	fmt.Printf("M:%v, R(%v)-P(%v)-E(%v)=%v, D:%v\n", mlen, rlen, plen, elen, rlen-plen-elen, dlen)
// 	return mlen, rlen, plen, elen, dlen
// }
