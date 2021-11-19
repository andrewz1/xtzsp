package xstorage

import (
	"net"
	"sync"
	"time"
)

const (
	storeTmo = time.Second * 10
	evictInt = time.Minute
)

type storage struct {
	sync.Mutex
	m map[skey]svalue
	t *time.Ticker
}

func NewStorage() *storage {
	s := &storage{
		m: make(map[skey]svalue),
		t: time.NewTicker(evictInt),
	}
	go s.evict()
	return s
}

func (s *storage) evictOnce(t int64) {
	s.Lock()
	for k, v := range s.m {
		if v.exp <= t {
			delete(s.m, k)
		}
	}
	s.Unlock()
}

func (s *storage) evict() {
	for tm := range s.t.C {
		s.evictOnce(tm.UnixNano())
	}
}

func (s *storage) Put(devIP, srcIP, dstIP net.IP, srcPort, dstPort uint16, seq uint32, srcMac net.HardwareAddr) {
	k := newKey(devIP, dstIP, dstPort, seq)
	v := newValue(srcMac, srcIP, srcPort)
	s.Lock()
	s.m[k] = v
	s.Unlock()
}

func (s *storage) Get(devIP, dstIP net.IP, dstPort uint16, seq uint32) (net.IP, uint16, net.HardwareAddr, bool) {
	k := newKey(devIP, dstIP, dstPort, seq)
	s.Lock()
	defer s.Unlock()
	if v, ok := s.m[k]; ok {
		delete(s.m, k)
		return v.srcIp, v.srcPort, v.srcMac, true
	}
	return nil, 0, nil, false
}
