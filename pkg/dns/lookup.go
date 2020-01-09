package dns

import (
	"context"
	"fmt"
	"github.com/bbs-team/Ndig/pkg/http"
	"github.com/bbs-team/Ndig/pkg/paint"
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

type lookupResult struct {
	dnsServer string
	country   string
	*result
}

func Do(domain string) *Lookup {
	rc := make(chan lookupResult)
	//ic := make(chan string)

	t := paint.New()
	t.SetHeader(table.Row{"Country", "DNS Server IP", "CNAME",  "A Records"})
	t.SetFooter(table.Row{"DOMAIN", domain})

	res := &Lookup{
		paint: t,
		done:  make(chan struct{}),
		idx:   len(countryMap),
	}

	go gatherResult(res, rc)

	// global public dns server 별 lookup
	for country, countryId := range countryMap {
		go do(domain, country, countryId, rc)
	}
	<-res.done

	return res
}

func do(domain, country, countryId string, queue chan lookupResult) {
	// http client 생성
	c := http.New()

	// GET public dns ip
	err := c.SetRequest(method, createURL(countryId)).Do()
	if err != nil {
		panic(err)
	}
	pDns := &publicDns{
		server: make([]dnsServer, 0),
	}

	err = c.UnmarshalJSON(&pDns.server)
	if err != nil {
		panic(err)
	}
	
	// dns server 선정
	pDns.designate(countryId)
	// dns client 생성
	d := New(pDns.ip)
	ctx, _ := context.WithTimeout(context.Background(), d.timeout)

	res, err := d.lookup(ctx, domain)
	if err != nil {
		e := fmt.Sprintf("%s/%s\n -> (%s)", country, pDns.ip, err)
		log.Println(e)
		queue <- lookupResult{
			dnsServer: pDns.ip,
			country:   country,
			result:    &result{
				A:     []net.IPAddr{
					{
						IP:   []byte{},
						Zone: "",
					},
				},
				CNAME: "",
			},
		}
	}

	data := lookupResult{
		dnsServer: pDns.ip,
		country:   country,
		result:    res,
	}

	queue <- data
}

func gatherResult(r *Lookup, c <- chan lookupResult) {
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
func (lr *lookupResult)pretty() string {
	ret := ""
	for i, v := range lr.A {
		if i == len(lr.A) -1 {
			ret += strings.ReplaceAll(strings.ReplaceAll(fmt.Sprint(v), "{", ""), "}", "")
			break
		}
	}

	return ret
}

func createURL(countryId string) string {
	// use, usw 는 us로 변경
	queryId := countryId

	if strings.Contains(countryId, "us") {
		queryId = "us"
	}

	return endpoint + queryId + dataFormat
}