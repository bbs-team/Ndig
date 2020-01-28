package dns

import (
	"context"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Client struct {
	udp *net.Resolver
	tcp *net.Resolver
	timeout  time.Duration
}

type result struct {
	A     []net.IPAddr
	CNAME string
}

func init() {
	log.SetOutput(os.Stderr)
}

func New(ip string) *Client {
	dc :=&Client{
		timeout: time.Second * 5,
	}

	dc.udp = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, n, a string) (conn net.Conn, e error) {
			d := &net.Dialer{
				Timeout:dc.timeout,
			}
			return d.DialContext(ctx, "udp", ip+":53")
		},
	}

	dc.tcp = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, n, a string) (conn net.Conn, e error) {
			d := &net.Dialer{
				Timeout:dc.timeout+5,
			}
			return d.DialContext(ctx, "tcp", ip+":53")
		},
	}

	return dc
}

func (c *Client) lookup(ctx context.Context, domain string) (*result, error) {
	defer ctx.Done()

	r, err := c.lookupUDP(ctx, domain)
	if err != nil {
		r, err = c.lookupTCP(ctx, domain)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (c *Client) lookupUDP(ctx context.Context, domain string) (*result, error) {
	// cname resolve
	cn, err := c.udp.LookupCNAME(ctx, domain)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(cn[:len(cn) - 1]) == strings.ToLower(domain) {
		cn = "-"
	}

	// ip resolve
	a, err := c.udp.LookupIPAddr(ctx, domain)
	if err != nil {
		return nil, err
	}

	r := &result{
		A: a,
		CNAME: cn,
	}

	return r, nil
}

func (c *Client) lookupTCP(ctx context.Context, domain string) (*result, error) {
	// cname resolve
	cn, err := c.tcp.LookupCNAME(ctx, domain)
	if err != nil {
		return nil, err
	}
	if cn[:len(cn) - 1] == domain {
		cn = "-"
	}

	// ip resolve
	a, err := c.tcp.LookupIPAddr(ctx, domain)
	if err != nil {
		return nil, err
	}

	r := &result{
		A: a,
		CNAME: cn,
	}

	return r, nil
}

func (c *Client) setTimeout(d time.Duration)  {
	c.timeout = d
}