package xstorage

import (
	"net"
)

type skey struct {
	devIP   string
	dstIP   string
	dstPort uint16
	seq     uint32
}

func newKey(devIP, dstIP net.IP, dstPort uint16, seq uint32) skey {
	return skey{
		devIP:   string(devIP),
		dstIP:   string(dstIP),
		dstPort: dstPort,
		seq:     seq,
	}
}
