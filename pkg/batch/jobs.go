package batch

import (
	"encoding/json"
	"github.com/bbs-team/ndig/pkg/dns"
	"log"
	"os"
	"sync"
	"time"
)

func UpdatePublicDnsFunc()  {
	// var init
	date := time.Now().Format("20060102")
	metaPath := "meta/"
	fileName := "_public-dns.json"
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	wg.Add(len(dns.CountryMap))
	s := &dns.Public{
		Servers:make(map[string]dns.Server, 0),
	}

	// Get public dns server information each countries
	for sName, fName := range dns.CountryMap {
		go func(sn, fn string) {
			mu.Lock()
			srv := dns.LoadPublicDns(sn)
			s.Servers[sn] = dns.Server{
				IP:      srv.Ip,
				Country: fn,
			}
			mu.Unlock()
			wg.Done()
		}(sName, fName)
	}
	wg.Wait()

	// Marshal data for write file
	data, err := json.MarshalIndent(&s, "", "\t")
	if err != nil {
		log.Fatal("updatePublicDns().json.Unmarshal():", err)
	}

	err := os.MkdirAll("meta", 0666)
	if err != nil {
		log.Fatal("updatePublicDns().os.MkdirAll():", err)
	}

	// create and write daily file
	f, err := os.Create(metaPath + date + fileName)
	if err != nil {
		log.Fatal("updatePublicDns().os.Open():", err)
	}

	_, err = f.WriteAt(data, 0)
	if err != nil {
		log.Fatal("updatePublicDns().os.WriteAt:", err)
	}

	defer f.Close()
}
