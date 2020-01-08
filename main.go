package main

import (
	"fmt"
	"log"
	"os"
	"query-gpop/pkg/dns"
)

func main()  {
	if len(os.Args) < 2 {
		printUsage()
		return
	}
	log.SetOutput(os.Stderr)
	domain := os.Args[1]

	//terminal.Clear()
	fmt.Printf("QUERY %s\nGET Public DNS Server IPs", domain)

	//terminal.Clear()
	dns.Do(domain).Print()
}

func printUsage()  {
	fmt.Println("usage: ./query-gpop [Domain]")
}