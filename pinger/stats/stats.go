package stats

import (
	"net"
	"time"

	"github.com/antonovegorv/simple-pinger/pinger/packet"
)

type Stats struct {
	DestIP             net.IP
	ReceivedPackets    []*packet.Packet
	PacketsTransmitted int
}

func (s *Stats) GetMinRTT() time.Duration {
	if len(s.ReceivedPackets) == 0 {
		return 0
	}

	minRTT := s.ReceivedPackets[0].RTT
	for _, p := range s.ReceivedPackets {
		if p.RTT < minRTT {
			minRTT = p.RTT
		}
	}

	return minRTT
}

func (s *Stats) GetMaxRTT() time.Duration {
	if len(s.ReceivedPackets) == 0 {
		return 0
	}

	maxRTT := s.ReceivedPackets[0].RTT
	for _, p := range s.ReceivedPackets {
		if p.RTT > maxRTT {
			maxRTT = p.RTT
		}
	}

	return maxRTT
}

func (s *Stats) GetAvgRTT() time.Duration {
	if len(s.ReceivedPackets) == 0 {
		return 0
	}

	var totalRTT time.Duration
	for _, p := range s.ReceivedPackets {
		totalRTT += p.RTT
	}

	return time.Duration(int(totalRTT) / len(s.ReceivedPackets))
}

func (s *Stats) GetPacketsLost() float64 {
	return ((float64(s.PacketsTransmitted) - float64(len(s.ReceivedPackets))) /
		float64(s.PacketsTransmitted)) * 100
}
