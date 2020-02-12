package dns

import (
	"context"
	"net"
	"strings"
)

type Result struct {
	A     []net.IPAddr
	CNAME string
}

func (c *Client) Lookup(ctx context.Context, domain string) (*Result, error) {
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

func (c *Client) lookupUDP(ctx context.Context, domain string) (*Result, error) {
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

	r := &Result{
		A: a,
		CNAME: cn,
	}

	return r, nil
}

func (c *Client) lookupTCP(ctx context.Context, domain string) (*Result, error) {
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

	r := &Result{
		A: a,
		CNAME: cn,
	}

	return r, nil
}