package main

import (
	"flag"
	"fmt"
	concordances "github.com/whosonfirst/go-whosonfirst-concordances"
	"os"
	"runtime"
	"strings"
)

func main() {

	var source = flag.String("source", "", "Where to look for files")
	var procs = flag.Int("processes", runtime.NumCPU()*2, "Number of concurrent processes to use")

	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	list := concordances.ListConcordances(*source)
	fmt.Println(strings.Join(list, ","))

	os.Exit(0)
}
