package main

import (
	"github.com/bbs-team/ndig/pkg/api"
	"os"
)

func main()  {
	s := api.NewServer(os.Getenv("PORT"))
	s.Start()
}
