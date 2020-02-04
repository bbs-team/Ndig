package dns

import (
	"context"
	"fmt"
	"github.com/bbs-team/ndig/pkg/http"
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
	defer func() {
		recv := recover()
		if recv != nil {
			log.Printf("do() error occurred: \n%s\n\n" , recv)
		}
	}()

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
		queue <- lookupResult{
			dnsServer: pDns.ip,
			country:   country,
			result:    &result{
				A:     []net.IPAddr{{},},
				CNAME: "",
			},
		}
		panic(e)
	}

	data := lookupResult{
		dnsServer: pDns.ip,
		country:   country,
		result:    res,
	}

	queue <- data
}

func gatherResult(r *Lookup, c <- chan lookupResult) {
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
func (lr *lookupResult)pretty() string {
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

func createURL(countryId string) string {
	// use, usw 는 us로 변경
	queryId := countryId

	if strings.Contains(countryId, "us") {
		queryId = "us"
	}

	return endpoint + queryId + dataFormat
}

func trimBrace(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "{", ""), "}", "")
}