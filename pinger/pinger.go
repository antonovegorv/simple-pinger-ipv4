package pinger

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/antonovegorv/simple-pinger/pinger/config"
	"github.com/antonovegorv/simple-pinger/pinger/packet"
	"github.com/antonovegorv/simple-pinger/pinger/stats"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const network = "ip4:icmp"
const address = "0.0.0.0"
const protocolICMP = 1

type Pinger struct {
	ctx        context.Context
	wg         *sync.WaitGroup
	errorsChan chan error
	config     *config.Config
	stats      stats.Stats
}

func New(ctx context.Context, wg *sync.WaitGroup, errorsChan chan error,
	config *config.Config) *Pinger {
	return &Pinger{
		ctx:        ctx,
		wg:         wg,
		errorsChan: errorsChan,
		config:     config,
	}
}

func (p *Pinger) Ping() {
	defer p.wg.Done()

	c, err := icmp.ListenPacket(network, address)
	if err != nil {
		p.errorsChan <- err
		return
	}
	defer c.Close()

	c.IPv4PacketConn().SetTTL(p.config.TTL)
	c.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true)

	ips, err := net.LookupIP(p.config.Hostname)
	if err != nil {
		p.errorsChan <- err
		return
	}

	for _, ip := range ips {
		if p.stats.DestIP = ip.To4(); p.stats.DestIP != nil {
			break
		}
	}

	if p.stats.DestIP == nil {
		p.errorsChan <- fmt.Errorf("no ipv4 for that host %v;", p.config.Hostname)
	}

	for i := 0; ; i++ {
		if i == p.config.Count && p.config.Count != 0 {
			break
		}

		select {
		case <-p.ctx.Done():
			return
		default:
			wm := icmp.Message{
				Type: ipv4.ICMPTypeEcho, Code: 0,
				Body: &icmp.Echo{
					ID: os.Getpid() & 0xffff, Seq: i + 1,
					Data: []byte(strings.Repeat("0", p.config.Size)),
				},
			}

			wb, err := wm.Marshal(nil)
			if err != nil {
				p.errorsChan <- err
				return
			}

			start := time.Now()

			if _, err := c.WriteTo(wb, &net.IPAddr{IP: p.stats.DestIP}); err != nil {
				p.errorsChan <- err
				return
			}

			p.stats.PacketsTransmitted++

			rb := make([]byte, 1500)

			var cm *ipv4.ControlMessage
			n, cm, peer, err := c.IPv4PacketConn().ReadFrom(rb)
			if err != nil {
				p.errorsChan <- err
				return
			}

			elapsed := time.Since(start)

			rm, err := icmp.ParseMessage(protocolICMP, rb[:n])
			if err != nil {
				p.errorsChan <- err
				return
			}

			switch rm.Type {
			case ipv4.ICMPTypeEchoReply:
				aPacket := packet.New(p.config.Size, peer, i+1, cm.TTL, elapsed)
				p.stats.ReceivedPackets = append(p.stats.ReceivedPackets, aPacket)
				aPacket.Log()
			default:
				fmt.Printf("got %+v; want echo reply", rm)
			}

			time.Sleep(time.Duration(p.config.Interval) * time.Second)
		}
	}

	p.errorsChan <- nil
}

func (p *Pinger) LogStats() {
	fmt.Printf("\n--- %v ping statistics ---\n", p.config.Hostname)
	fmt.Printf("%v packets transmitted, %v received, %v packet lost\n",
		p.stats.PacketsTransmitted, len(p.stats.ReceivedPackets),
		p.stats.GetPacketsLost())
	fmt.Printf("rtt min/avg/max = %v/%v/%v\n",
		p.stats.GetMinRTT(), p.stats.GetAvgRTT(), p.stats.GetMaxRTT())
}
