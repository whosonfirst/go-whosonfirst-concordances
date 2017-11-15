package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-concordances"
	"log"
	"os"
	"runtime"
	"strings"		
)

func main() {

	mode := flag.String("mode", "repo", "")
	var procs = flag.Int("processes", runtime.NumCPU()*2, "Number of concurrent processes to use")

	runtime.GOMAXPROCS(*procs)
	
	flag.Parse()

	sources := flag.Args()

	list, err := concordances.ListConcordances(*mode, sources...)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.Join(list, ","))
	os.Exit(0)
}
