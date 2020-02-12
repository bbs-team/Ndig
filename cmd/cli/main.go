package main

import (
	"fmt"
	"github.com/bbs-team/ndig/pkg/cli"
	"log"
	"os"
)

const appName = "ndig"
const version  = "0.1"

func main()  {
	if len(os.Args) < 2 {
		printUsage()
		return
	}
	log.SetOutput(os.Stderr)
	domain := os.Args[1]

	//terminal.Clear()
	fmt.Printf("%s v%s\nQUERY %s\n\n",appName, version, domain)

	//terminal.Clear()
	cli.Do(domain).Print()
}

func printUsage()  {
	fmt.Println("usage: ./ndig [Domain]")
}