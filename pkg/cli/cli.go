package cli

import (
	"context"
	"fmt"
	"github.com/bbs-team/ndig/pkg/dns"
	"github.com/bbs-team/ndig/pkg/paint"
	"github.com/jedib0t/go-pretty/table"
	"log"
	"net"
	"strings"
)

type Lookup struct {
	paint *paint.Paint
	idx int
	done chan struct{}
}

type LookupResult struct {
	dnsServer string
	country   string
	*dns.Result
}

func Do(domain string) *Lookup {

	rc := make(chan LookupResult)
	//ic := make(chan string)

	t := paint.New()
	t.SetHeader(table.Row{"Country", "DNS Server IP", "CNAME",  "A Records"})
	t.SetFooter(table.Row{"DOMAIN", domain})

	res := &Lookup{
		paint: t,
		done:  make(chan struct{}),
		idx:   len(dns.CountryMap),
	}

	go gatherResult(res, rc)

	// global public dns server 별 lookup
	for countryId, country := range dns.CountryMap {
		go do(domain, country, countryId, rc)
	}
	<-res.done
	close(rc)

	return res
}

func do(domain, country, countryId string, queue chan LookupResult) {
	defer func() {
		recv := recover()
		if recv != nil {
			log.Printf("do() error occurred: \n%s\n\n" , recv)
		}
	}()

	pDns := dns.LoadPublicDns(countryId)
	// dns client 생성
	d := dns.New(pDns.Ip)
	ctx, _ := context.WithTimeout(context.Background(), d.Timeout)

	res, err := d.Lookup(ctx, domain)
	if err != nil {
		e := fmt.Sprintf("%s/%s\n -> (%s)", country, pDns.Ip, err)
		queue <- LookupResult{
			dnsServer: pDns.Ip,
			country:   country,
			Result:    &dns.Result{
				A:     []net.IPAddr{{},},
				CNAME: "",
			},
		}
		panic(e)
	}

	data := LookupResult{
		dnsServer: pDns.Ip,
		country:   country,
		Result:    res,
	}

	queue <- data
}

func gatherResult(r *Lookup, c <- chan LookupResult) {
	defer func() {
		recv := recover()
		if recv != nil {
			log.Printf("do() error occurred: \n%s\n\n" , recv)
		}
	}()

	for {
		res := <- c
		r.paint.SetRow(table.Row{res.country, res.dnsServer, res.CNAME, res.pretty()})
		r.decreaseIdx()
		if r.idx == 0 {
			r.done <- struct{}{}
			break
		}
	}
}

func (l *Lookup)decreaseIdx()  {
	l.idx = l.idx - 1
}

func (l *Lookup) Print() {
	l.paint.Render()
}


// A 레코드를 보기 좋게..
func (lr *LookupResult)pretty() string {
	ret := ""


	for i, v := range lr.A {
		t := fmt.Sprint(v)
		if t == "{<nil> }" {
			return trimBrace(t)
		}

		if i == len(lr.A) -1 {
			ret += trimBrace(t)
			break
		}

		ret += trimBrace(t) + "\n"
	}
	return ret
}

func trimBrace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "{", ""), "}", "")
}