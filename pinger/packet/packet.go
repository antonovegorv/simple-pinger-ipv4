package packet

import (
	"fmt"
	"net"
	"time"
)

type Packet struct {
	Size int
	Peer net.Addr
	Seq  int
	TTL  int
	RTT  time.Duration
}

func New(size int, peer net.Addr, seq, ttl int, rtt time.Duration) *Packet {
	return &Packet{
		Size: size,
		Peer: peer,
		Seq:  seq,
		TTL:  ttl,
		RTT:  rtt,
	}
}

func (p *Packet) Log() {
	fmt.Printf("%v bytes from %v: icmp_seq=%v ttl=%v time=%v\n",
		p.Size, p.Peer.String(), p.Seq, p.TTL, p.RTT)
}
