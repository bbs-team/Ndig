package batch

import "github.com/robfig/cron/v3"

func Start()  {
	c := cron.New()
	_, err := c.AddFunc("@daily", UpdatePublicDnsFunc)
	if err != nil {
		panic(err)
	}

	c.Start()
}

