package xrecv

import (
	"fmt"
	"net"

	"github.com/andrewz1/xtzsp/xpacket"
)

const (
	listen = ":5514"
	bufLen = 4096
)

var (
	hdr    = []byte{1, 0, 0, 1, 1} // default header
	hdrS   = string(hdr)
	hdrLen = len(hdr)
)

func Init() error {
	var ua *net.UDPAddr
	var err error
	if ua, err = net.ResolveUDPAddr("udp", listen); err != nil {
		return err
	}
	var uc *net.UDPConn
	if uc, err = net.ListenUDP("udp", ua); err != nil {
		return err
	}
	go receiver(uc)
	return nil
}

func hdrIsOK(buf []byte) bool {
	if len(buf) < hdrLen {
		return false
	}
	return string(buf[:hdrLen]) == hdrS
}

func receiver(uc *net.UDPConn) {
	var buf [bufLen]byte
	var ua *net.UDPAddr
	var n int
	var err error
	for {
		if n, ua, err = uc.ReadFromUDP(buf[:]); err != nil || !hdrIsOK(buf[:n]) {
			continue
		}
		go process(append(net.IP{}, ua.IP...), append([]byte{}, buf[hdrLen:n]...))
	}
}

func process(devIP net.IP, buf []byte) {
	if p := xpacket.NewPacket(devIP, buf); p != nil {
		fmt.Println(p)
	}
}
