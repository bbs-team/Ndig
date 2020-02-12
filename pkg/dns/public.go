package dns

import (
	"github.com/bbs-team/ndig/pkg/http"
	"strings"
)

const endpoint  = "https://public-dns.info/nameserver/"
const method  = "GET"
const dataFormat = ".json"

/*var CountryMap = map[string]string{
	"Korea":              "kr",
	"Japan":              "jp",
	"China":              "cn",
	"Hong Kong":          "hk",
	"Singapore":          "sg",
	"Thailand":           "th",
	"Taiwan":             "tw",
	"Vietnam":            "vn",
	"USW(San Francisco)": "usw",
	"USE(New York)":      "use",
	"Germany":            "de",
	"France":             "fr",
	"United Kingdom":     "gb",
	"EU(Argentina)":      "ar",
	"Indonesia":          "id",
}*/

var CountryMap = map[string]string{
	"kr":  "Korea",
	"jp":  "Japan",
	"cn":  "China",
	"hk":  "Hong Kong",
	"sg":  "Singapore",
	"th":  "Thailand",
	"tw":  "Taiwan",
	"vn":  "Vietnam",
	"usw": "USW(San Francisco)",
	"use": "USE(New York)",
	"de":  "Germany",
	"fr":  "France",
	"gb":  "United Kingdom",
	"ar":  "EU(Argentina)",
	"id":  "Indonesia",
}

type PublicApiResponse struct {
	CountryId   string  `json:"country_id"`
	City		string  `json:"city"`
	Ip          string  `json:"ip"`
	Reliability float32 `json:"reliability"`
	Name        string  `json:"name"`
}

func LoadPublicDns(countryId string) *PublicApiResponse {
	// http client 생성
	c := http.New()

	// GET public dns ip
	err := c.SetRequest(method, createURL(countryId)).Do()
	if err != nil {
		panic(err)
	}
	dnsServers := make([]PublicApiResponse, 0)

	err = c.UnmarshalJSON(&dnsServers)
	if err != nil {
		panic(err)
	}

	// dns server 선정
	return designate(dnsServers, countryId)
}

func designate(servers []PublicApiResponse, countryId string) *PublicApiResponse {
	dnsServer := &PublicApiResponse{}

	for _, v := range servers {

		// google dns는 패스
		if  v.Ip == "8.8.8.8" {
			continue
		}

		/*if countryId == "id" {
			ip = "203.189.88.233"
			break
		}*/

		// usw는 샌프란시스코
		if countryId == "usw" &&
			(v.City == "San Francisco" || v.Reliability > 1) {
			dnsServer = &v
			break
		}
		// use는 뉴욕
		if countryId == "use" &&
			(v.City == "New York" || v.Reliability > 1) {
			dnsServer = &v
			break
		}

		// 기타 다른 국가는 도시 무관
		if (countryId != "use" && countryId != "usw") ||
			v.Reliability > 1 {
			dnsServer = &v
			break
		}
	}

	return dnsServer
}

func createURL(countryId string) string {
	// use, usw 는 us로 변경
	queryId := countryId

	if strings.Contains(countryId, "us") {
		queryId = "us"
	}

	return endpoint + queryId + dataFormat
}