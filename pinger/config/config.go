package config

type Config struct {
	Hostname string
	Interval int
	Count    int
	TTL      int
	Size     int
}

func New(hostname string, interval, count, ttl, size int) *Config {
	return &Config{
		Hostname: hostname,
		Interval: interval,
		Count:    count,
		TTL:      ttl,
		Size:     size,
	}
}
