package main

import (
	"flag"
	concordances "github.com/whosonfirst/go-whosonfirst-concordances"
	"os"
	"runtime"
)

func main() {

	var source = flag.String("source", "https://whosonfirst.mapzen.com/data/", "Where to look for files")
	var procs = flag.Int("processes", 100, "Number of concurrent processes to use")

	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	concordances.WriteConcordances(*source, os.Stdout)

	os.Exit(0)
}
