package util

import (
	"encoding/json"
	"github.com/bbs-team/ndig/pkg/dns"
	"io/ioutil"
	"time"
)

func LoadDnsFile() (*dns.Public, error) {
	date := time.Now().Format("20060102")
	path := "meta/"
	file := "_public-dns.json"

	b, err := ioutil.ReadFile(path + date + file)
	if err != nil {
		return nil, err
	}

	publicDns := &dns.Public{
		Servers:make(map[string]dns.Server, 0),
	}
	err = json.Unmarshal(b, publicDns)
	if err != nil {
		return nil, err
	}

	return publicDns, nil
}
