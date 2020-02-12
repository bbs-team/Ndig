package api

import (
	"context"
	"github.com/bbs-team/ndig/pkg/dns"
	"github.com/bbs-team/ndig/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

func (s *server) help() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := HelpResponse{
			Functions: make([]function, 0),
		}
		for _, v := range s.Routes() {
			resp.Functions = append(resp.Functions, function{
				Path:   v.Path,
				Method: v.Method,
			})
		}
		response(c, resp, http.StatusOK, "")
		return
	}
}

func countries() gin.HandlerFunc {
	return func(c *gin.Context) {
		response(c, dns.CountryMap, http.StatusOK, "ok")
	}
}

func queryDns() gin.HandlerFunc {
	return func(c *gin.Context) {
		domain := c.Param("domain")
		country := c.Query("c")

		publicDns, err := util.LoadDnsFile()
		if err != nil {
			status := http.StatusInternalServerError
			msg := "Public DNS list loading is failed"
			response(c, nil, status, msg)
			return
		}

		resp := make([]QueryResponse, 0)
		wg := &sync.WaitGroup{}
		mu := &sync.Mutex{}
		// 특정 국가 쿼리 요청
		if country != "" {
			switch publicDns.Servers[country].Country {
			case "":
				status := http.StatusNotFound
				msg := country + " is not support country code"
				response(c, nil, status, msg)
				return
			default:
				respData := queryProcess(domain, publicDns.Servers[country])
				resp = append(resp, *respData)
				response(c, resp, http.StatusOK, "ok")
				return
			}
		}

		// 전체 국가 쿼리 요청
		for _, v := range publicDns.Servers {
			wg.Add(1)
			go func(ds dns.Server) {
				respData := queryProcess(domain, ds)
				mu.Lock()
				resp = append(resp, *respData)
				mu.Unlock()
				wg.Done()
			}(v)
		}

		wg.Wait()

		response(c, resp, http.StatusOK, "ok")
		return
	}
}

func queryProcess(domain string, srv dns.Server) *QueryResponse {
	dnsClient := dns.New(srv.IP)
	ctx, _ := context.WithTimeout(context.Background(), dnsClient.Timeout)
	lookupRes, err := dnsClient.Lookup(ctx, domain)
	qr := &queryResult{}
	if lookupRes != nil {
		qr.Cname = lookupRes.CNAME
		for _, v := range lookupRes.A {
			qr.IP = append(qr.IP, v.IP.String())
		}
	}
	respData := &QueryResponse{
		Country:     srv.Country,
		DnsServerIP: srv.IP,
		Results:     qr,
	}
	if err != nil {
		respData.QueryError = err.Error()
		respData.Results = nil
	}

	return respData
}