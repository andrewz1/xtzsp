package xstorage

import (
	"net"
	"time"
)

type svalue struct {
	srcMac  net.HardwareAddr
	srcIp   net.IP
	srcPort uint16
	exp     int64
}

func newValue(srcMac net.HardwareAddr, srcIp net.IP, srcPort uint16) svalue {
	return svalue{
		srcMac:  append(net.HardwareAddr{}, srcMac...),
		srcIp:   append(net.IP{}, srcIp...),
		srcPort: srcPort,
		exp:     time.Now().Add(storeTmo).UnixNano(),
	}
}
