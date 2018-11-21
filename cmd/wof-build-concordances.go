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

	mode := flag.String("mode", "repo", "")
	out := flag.String("outfile", "", "")

	var procs = flag.Int("processes", runtime.NumCPU()*2, "Number of concurrent processes to use")

	flag.Parse()

	runtime.GOMAXPROCS(*procs)

	sources := flag.Args()

	if *mode == "repo" && len(sources) == 0 {

		cwd, err := os.Getwd()

		if err != nil {
			log.Fatal(err)
		}

		root := cwd

		abs_repo, err := filepath.Abs(root)

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

		r_opts := repo.DefaultFilenameOptions()
		r, err := repo.NewDataRepoFromPath(abs_repo, r_opts)

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

		*out = filepath.Join(abs_meta, fname)

		sources = append(sources, abs_repo)
	}

	if *out == "" {
		log.Fatal("missing outfile")
	}

	fh, err := atomicfile.New(*out, os.FileMode(0644))

	if err != nil {
		log.Fatal(err)
	}

	err = concordances.WriteConcordances(fh, *mode, sources...)

	if err != nil {
		fh.Abort()
		log.Fatal(err)
	}

	fh.Close()
	os.Exit(0)
}
