package xpacket

import (
	"net"

	"github.com/andrewz1/xtzsp/xstorage"
)

type storage interface {
	Put(devIP, srcIP, dstIP net.IP, srcPort, dstPort uint16, seq uint32, srcMac net.HardwareAddr)
	Get(devIP, dstIP net.IP, dstPort uint16, seq uint32) (srcIP net.IP, srcPort uint16, srcMac net.HardwareAddr, ok bool)
}

var s storage

func init() {
	s = xstorage.NewStorage()
}
