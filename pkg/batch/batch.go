package batch

import (
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

var logWriter *log.Logger

func init()  {
	logPath := "log/"
	logFile := "batch_job.log"

	err := os.MkdirAll("log", 0777)
	if err != nil {
		log.Fatal(err)
	}

	// Logging to a file.
	f, err := os.Create(logPath + logFile)
	if err != nil {
		log.Fatal(err)
	}
	logWriter = log.New(f, "cron: ", log.LstdFlags)
}

func Start()  {
	c := cron.New()

	_, err := c.AddFunc("@daily", UpdatePublicDnsFunc)
	if err != nil {
		panic(err)
	}

	c.Start()
}

