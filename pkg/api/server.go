package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"time"
)

type Server interface {
	Start()
}

type server struct {
	*gin.Engine
	port string
}

func NewServer(port string) Server {
	gin.DisableConsoleColor()
	logPath := "log/"
	logFile := "_access.log"
	date := time.Now().Format("20060102")

	err := os.MkdirAll("log", 0777)
	if err != nil {
		log.Fatal(err)
	}

	// Logging to a file.
	f, err := os.Create(logPath + date + logFile)
	if err != nil {
		log.Fatal(err)
	}

	gin.DefaultWriter = io.MultiWriter(f)

	e := gin.New()
	e.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] - %s %s \"%s %s %s %s %d\" %d %d %d %s %s\" %s\n",
			param.TimeStamp.Format(time.RFC1123),
			param.ClientIP,
			param.Request.Host,
			param.Method,
			param.Request.URL.Path,
			param.Request.URL.RawQuery,
			param.Request.Proto,
			param.StatusCode,
			param.Request.ContentLength,
			param.BodySize,
			param.Latency.Milliseconds(),
			param.Request.UserAgent(),
			param.Request.Referer(),
			param.ErrorMessage,
		)
	}),
		gin.Recovery())

	return &server{
		Engine: e,
		port:   ":" + port,
	}
}

func (s *server)Start() {
	// Disable Console Color

	s.controller()
	s.Run(s.port)
}

func (s *server)controller()  {
	// apis
	s.NoRoute(s.help())
	dnsApi := s.Group("/dns")
	{
		dnsApi.GET("/query/:domain", queryDns())
		dnsApi.GET("/test", ResponseHtml())
		dnsApi.GET("/countries", countries())
	}
}