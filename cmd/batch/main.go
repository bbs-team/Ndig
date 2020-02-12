package main

import "github.com/bbs-team/ndig/pkg/batch"

func main()  {
	var done chan struct{}
	batch.Start()
	<-done
}
