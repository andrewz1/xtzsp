package xpacket

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/google/gopacket/layers"
)

type Packet struct {
	devIP   net.IP           // device IP
	srcMac  net.HardwareAddr // host src mac
	srcIP   net.IP           // orig host IP
	extIP   net.IP           // nat ext IP
	dstIP   net.IP           // dst IP
	srcPort uint16           // orig host port
	extPort uint16           // nat ext port
	dstPort uint16           // dst port
}

func NewPacket(devIP net.IP, buf []byte) *Packet {
	var eth layers.Ethernet
	var fb fback
	var err error
	if err = eth.DecodeFromBytes(buf, &fb); err != nil || fb.isTruncated() {
		return nil
	}
	var nbuf []byte
	switch eth.EthernetType {
	case layers.EthernetTypeLLC, layers.EthernetTypeIPv4:
		nbuf = eth.Payload
	case layers.EthernetTypeDot1Q:
		var dotq layers.Dot1Q
		if err = dotq.DecodeFromBytes(eth.Payload, &fb); err != nil || fb.isTruncated() {
			return nil
		}
		nbuf = dotq.Payload
	default: // unknown packet
		return nil
	}
	var ip4 layers.IPv4
	if err = ip4.DecodeFromBytes(nbuf, &fb); err != nil || fb.isTruncated() {
		return nil
	}
	if ip4.Protocol != layers.IPProtocolTCP {
		return nil
	}
	var tcp layers.TCP
	if err = tcp.DecodeFromBytes(ip4.Payload, &fb); err != nil {
		return nil
	}
	// first packet = SYN, second packet = SYN+ACK
	if !tcp.SYN {
		return nil
	}
	if !tcp.ACK { // first packet with original host mac,ip,port
		s.Put(devIP, ip4.SrcIP, ip4.DstIP, uint16(tcp.SrcPort), uint16(tcp.DstPort), tcp.Seq+1, eth.SrcMAC)
		return nil
	}
	srcIP, srcPort, srcMac, ok := s.Get(devIP, ip4.SrcIP, uint16(tcp.SrcPort), tcp.Ack)
	if !ok {
		return nil // packet not found
	}
	return &Packet{
		devIP:   devIP,
		srcMac:  srcMac,
		srcIP:   srcIP,
		extIP:   append(net.IP{}, ip4.DstIP...),
		dstIP:   append(net.IP{}, ip4.SrcIP...),
		srcPort: srcPort,
		extPort: uint16(tcp.DstPort),
		dstPort: uint16(tcp.SrcPort),
	}
}

func writeMac(w io.Writer, mac net.HardwareAddr) {
	if len(mac) != 6 {
		fmt.Fprint(w, "<invalid mac>")
		return
	}
	fmt.Fprintf(w, "%02x-%02x-%02x-%02x-%02x-%02x", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

func (p *Packet) String() string {
	if p == nil {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "[%s] ", p.devIP)
	writeMac(&b, p.srcMac)
	fmt.Fprintf(&b, " %s:%d", p.srcIP, p.srcPort)
	fmt.Fprintf(&b, " %s:%d", p.extIP, p.extPort)
	fmt.Fprintf(&b, " %s:%d", p.dstIP, p.dstPort)
	return b.String()
}
