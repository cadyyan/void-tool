package configuration

import (
	"net"
	"time"
)

type HTTPConfiguration struct {
	Host            string        `default:"0.0.0.0"`
	Port            string        `default:"8080"`
	ShutdownTimeout time.Duration `default:"30s"`
	RequestTimeout  time.Duration `default:"1m"`
}

func (config HTTPConfiguration) BindAddress() string {
	return net.JoinHostPort(config.Host, config.Port)
}
