package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)


type Response struct {
	Data  interface{} `json:"data"`
	Error ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func response(c *gin.Context, data interface{}, status int, message string)  {
	resp := &Response{
		Data: data,
		Error:ErrorResponse{
			Code:    status,
			Message: message,
		},
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Render(200, render.IndentedJSON{Data:resp})
}

type HelpResponse struct {
	Functions []function `json:"functions"`
}

type function struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type QueryResponse struct {
	Country     string       `json:"country"`
	DnsServerIP string       `json:"dnsServerIp"`
	Results     *queryResult `json:"results"`
	QueryError  string       `json:"queryError"`
}

type queryResult struct {
	Cname string   `json:"cname"`
	IP    []string `json:"a"`
}
