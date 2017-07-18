package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-concordances"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {

	source := flag.String("source", "", "Where to look for files (This flag is deprecated and you should use 'repo' instead.)")
	repo := flag.String("repo", "", "Where to read data (to create metafiles) from. If empty then the code will assume the current working directory.")
	var procs = flag.Int("processes", runtime.NumCPU()*2, "Number of concurrent processes to use")

	flag.Parse()

	if *source != "" && *repo == "" {

		log.Println("The -source flag is officially deprecated and you should use the -repo flag instead.")
		*repo = filepath.Dir(*source)
	}

	// sudo put me in a helper function

	if *repo == "" {

		cwd, err := os.Getwd()

		if err != nil {
			log.Fatal(err)
		}

		*repo = cwd
	}

	abs_repo, err := filepath.Abs(*repo)

	if err != nil {
		log.Fatal(err)
	}

	info, err := os.Stat(abs_repo)

	if err != nil {
		log.Fatal(err)
	}

	if !info.IsDir() {
		log.Fatal(fmt.Sprintf("Invalid repo directory (%s)", abs_repo))
	}

	abs_data := filepath.Join(abs_repo, "data")

	info, err = os.Stat(abs_data)

	if err != nil {
		log.Fatal(err)
	}

	if !info.IsDir() {
		log.Fatal(fmt.Sprintf("Invalid data directory (%s)", abs_data))
	}

	// end of sudo put me in a helper function

	runtime.GOMAXPROCS(*procs)

	list, err := concordances.ListConcordances(abs_data)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.Join(list, ","))
	os.Exit(0)
}
