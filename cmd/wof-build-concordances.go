package main

import (
	"flag"
	"fmt"
	"github.com/facebookgo/atomicfile"
	"github.com/whosonfirst/go-whosonfirst-concordances"
	"github.com/whosonfirst/go-whosonfirst-repo"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {

	root := flag.String("repo", "", "Where to read data (to create metafiles) from. If empty then the code will assume the current working directory.")
	out := flag.String("outfile", "", "Where to store metafiles. If empty then assume metafile are created in a child folder of 'repo' called 'meta'.")

	var procs = flag.Int("processes", runtime.NumCPU()*2, "Number of concurrent processes to use")

	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	// sudo put me in a helper function

	if *root == "" {

		cwd, err := os.Getwd()

		if err != nil {
			log.Fatal(err)
		}

		*root = cwd
	}

	abs_repo, err := filepath.Abs(*root)

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

	r, err := repo.NewDataRepoFromPath(abs_repo)

	if err != nil {
		log.Fatal(err)
	}

	var abs_meta string

	if *out == "" {
		abs_meta = filepath.Join(abs_repo, "meta")
	} else {
		abs_meta = *out
	}

	info, err = os.Stat(abs_meta)

	if err != nil {
		log.Fatal(err)
	}

	if !info.IsDir() {
		log.Fatal(fmt.Sprintf("Invalid meta directory (%s)", abs_meta))
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

	// THIS IS A HACK in advance of
	// a) removing pre-generated concordances from repos
	// b) updating the WriteConcordances logic to generate per-placetype
	//    concordances files (below)
	// (20170726/thisisaaronland)

	opts := repo.DefaultFilenameOptions()
	fname := r.ConcordancesFilename(opts)

	fname = strings.Replace(fname, "-all-", "-", -1)

	outfile := filepath.Join(abs_meta, fname)

	fh, err := atomicfile.New(outfile, os.FileMode(0644))

	if err != nil {
		log.Fatal(err)
	}

	err = concordances.WriteConcordances(abs_data, fh)

	if err != nil {
		fh.Abort()
		log.Fatal(err)
	}

	fh.Close()
	os.Exit(0)
}
