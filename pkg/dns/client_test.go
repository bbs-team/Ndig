package dns

import (
	"context"
	"fmt"
	"testing"
)

func TestLookup(t *testing.T) {
	ip := "203.189.88.233"
	domain := "resources-rmcnmv.pstatic.net"
	dnsClient := New(ip)
	if dnsClient == nil {
		t.Failed()
	}
	ctx := context.Background()
	res, err := dnsClient.lookup(ctx, domain)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}