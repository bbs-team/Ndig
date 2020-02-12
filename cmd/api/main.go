package main

import "github.com/bbs-team/ndig/pkg/api"

func main()  {
	s := api.NewServer("8080")
	s.Start()
}
