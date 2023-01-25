package types

import (
	"fmt"
	"net"
)

type GlobalOptions struct {
	Verbose bool `opts:"-v,--verbose" desc:"verbose logging"`
}

type ServerOptions struct {
	Port     string `opts:"-p,--port" desc:"server port" default:"8000"`
	Hostname string `opts:"-h,--hostname" desc:"server hostname" default:"localhost"`
}

type EchoOptions struct {
	Port string `opts:"-p,--port" desc:"server port" default:"7000"`
}

type ClientOptions struct {
	Target            string            `json:"target"      opts:"[0]"                desc:"target to proxy"`
	Name              string            `json:"name"        opts:"-n,--name"          desc:"name of the tunnel"`
	ServerHost        string            `json:"server-host" opts:"-s,--server-host"   desc:"server hostname"                  default:"localhost"`
	ServerPort        string            `json:"server-port" opts:"-p,--server-port"   desc:"server port"                      default:"8000"`
	Insecure          bool              `json:"insecure"    opts:"-k,--insecure"      desc:"use insecure HTTP and WebSockets" default:"true"`
	AllowedIPs        []string          `json:"allowed-ips" opts:"-a,--allow-ip"      desc:"IP CIDR ranges to allow"          default:"0.0.0.0/0,::/0"`
	ReconnectAttempts int               `json:"-"           opts:"-r,--max-reconnect" desc:"max reconnect attempts"           default:"5"`
	Headers           map[string]string `json:"headers"     opts:"-h,--header"        desc:"headers to add to requests"`
}

func (c ClientOptions) Origin() string {
	return c.SchemeHTTP() + "://" + c.ServerHost
}

func (c ClientOptions) URL() string {
	return c.SchemeWS() + "://" + c.ServerHost + ":" + c.ServerPort + "/register?name=" + c.Name
}

func (c ClientOptions) SchemeHTTP() string {
	if c.Insecure {
		return "http"
	}
	return "https"
}

func (c ClientOptions) SchemeWS() string {
	if c.Insecure {
		return "ws"
	}
	return "wss"
}

func (c ClientOptions) Valid() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.Target == "" {
		return fmt.Errorf("target is required")
	}
	for _, ip := range c.AllowedIPs {
		if _, _, err := net.ParseCIDR(ip); err != nil {
			return fmt.Errorf("invalid IP CIDR range specified: %s", ip)
		}
	}
	return nil
}
