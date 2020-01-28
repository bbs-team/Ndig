package dns

const endpoint  = "https://public-dns.info/nameserver/"
const method  = "GET"
const dataFormat = ".json"

var countryMap = map[string]string{
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
}

type publicDns struct {
	server []dnsServer
	ip string
}

type dnsServer struct {
	CountryId   string  `json:"country_id"`
	City		string  `json:"city"`
	Ip          string  `json:"ip"`
	Reliability float32 `json:"reliability"`
	Name        string  `json:"name"`
}

func (p *publicDns)designate(countryId string)  {
	ip := ""
	for _, v := range p.server {

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
			ip = v.Ip
			break
		}
		// use는 뉴욕
		if countryId == "use" &&
			(v.City == "New York" || v.Reliability > 1) {
			ip = v.Ip
			break
		}

		// 기타 다른 국가는 도시 무관
		if (countryId != "use" && countryId != "usw") ||
			v.Reliability > 1 {
			ip = v.Ip
			break
		}
	}

	p.ip = ip
}