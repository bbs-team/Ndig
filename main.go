package main

import (
	"fmt"
	"github.com/bbs-team/Ndig/pkg/dns"
	"log"
	"os"
)

func main()  {
	if len(os.Args) < 2 {
		printUsage()
		return
	}
	log.SetOutput(os.Stderr)
	domain := os.Args[1]

	//terminal.Clear()
	fmt.Printf("QUERY %s\nGET Public DNS Server IPs\n\n", domain)

	//terminal.Clear()
	dns.Do(domain).Print()
}

func printUsage()  {
	fmt.Println("usage: ./query-gpop [Domain]")
}