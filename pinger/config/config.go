package config

type Config struct {
	Hostname   string
	Timeout    int
	Interval   int
	Count      int
	TTL        int
	PacketSize int
}

func New(hostname string, timeout, interval, count, ttl, packetSize int) *Config {
	return &Config{
		Hostname:   hostname,
		Timeout:    timeout,
		Interval:   interval,
		Count:      count,
		TTL:        ttl,
		PacketSize: packetSize,
	}
}
