package dns

import (
	"context"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	udp *net.Resolver
	tcp *net.Resolver
	Timeout  time.Duration
}

func init() {
	log.SetOutput(os.Stderr)
}

func New(ip string) *Client {
	dc :=&Client{
		Timeout: time.Second * 5,
	}

	dc.udp = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, n, a string) (conn net.Conn, e error) {
			d := &net.Dialer{
				Timeout:dc.Timeout,
			}
			return d.DialContext(ctx, "udp", ip+":53")
		},
	}

	dc.tcp = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, n, a string) (conn net.Conn, e error) {
			d := &net.Dialer{
				Timeout:dc.Timeout+5,
			}
			return d.DialContext(ctx, "tcp", ip+":53")
		},
	}

	return dc
}

func (c *Client) setTimeout(d time.Duration)  {
	c.Timeout = d
}