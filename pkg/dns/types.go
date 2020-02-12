package dns

type Public struct {
	Servers map[string]Server `json:"servers"`
}

type Server struct {
	IP        string `json:"ip"`
	Country  string `json:"country"`
}