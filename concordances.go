package concordances

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-index/utils"
	"io"
	"io/ioutil"
	_ "log"
	"os"
	"sort"
	"strings"
	"sync"
)

func ListConcordances(mode string, sources ...string) ([]string, error) {

	tmp := make(map[string]int)
	concordances := make([]string, 0)

	mu := sync.Mutex{}

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		f, err := load_feature(fh, ctx)

		if err != nil {
			return err
		}

		if f == nil {
			return nil
		}

		c, err := whosonfirst.Concordances(f)

		if err != nil {
			return err
		}

		mu.Lock()

		for src, _ := range c {
			tmp[src] += 1
		}

		mu.Unlock()
		return nil
	}

	idx, err := index.NewIndexer(mode, cb)

	if err != nil {
		return concordances, err
	}

	err = idx.IndexPaths(sources)

	if err != nil {
		return concordances, err
	}

	for name, _ := range tmp {
		concordances = append(concordances, name)
	}

	return concordances, nil
}

func WriteConcordances(out io.Writer, mode string, sources ...string) error {

	// we do a first pass over all the files to build a keys dictionary
	// mapping (concordance) source and count and for every non-zero length
	// concordance dictionary we serialize the data to a tmpfile which we
	// then re-read and dump as a CSV file using the keys dictionary to
	// generate the CSV header (20171114/thisisaaronland)

	tmpfile, err := ioutil.TempFile("", "concordances")

	if err != nil {
		return err
	}

	defer os.Remove(tmpfile.Name())

	keys := make(map[string]int)

	mu := sync.Mutex{}

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		f, err := load_feature(fh, ctx)

		if err != nil {
			return err
		}

		if f == nil {
			return nil
		}

		c, err := whosonfirst.Concordances(f)

		if err != nil {
			return err
		}

		mu.Lock()
		defer mu.Unlock()

		added := 0

		for k, _ := range c {

			_, ok := keys[k]

			if !ok {
				keys[k] = 0
			}

			keys[k] += 1
			added += 1
		}

		if added == 0 {
			return nil
		}

		// this is a valid concordance so ensure that
		// the current wof:id is included

		// remember the interface for the geojson.Feature
		// Id() method is to return a string and not a
		// WOF-ish int64 (which would confuse the type
		// definition for whosonfirst.WOFConcordances

		c["wof:id"] = f.Id()

		_, ok := keys["wof:id"]

		if !ok {
			keys["wof:id"] = 0
		}

		keys["wof:id"] += 1

		enc_c, err := json.Marshal(c)

		if err != nil {
			return nil
		}

		tmpfile.Write(enc_c)
		tmpfile.Write([]byte("\n"))

		return nil
	}

	// here is where we plow through all the data

	idx, err := index.NewIndexer(mode, cb)

	if err != nil {
		return err
	}

	err = idx.IndexPaths(sources)

	if err != nil {
		return err
	}

	// here is where we rewind the tempfile and generate
	// a (CSV) header

	tmpfile.Seek(0, 0)

	header := make([]string, 0)

	for k, _ := range keys {
		header = append(header, k)
	}

	sort.Sort(sort.StringSlice(header))

	writer := csv.NewWriter(out)
	writer.Write(header)
	writer.Flush()

	// here is where we re-read the tempfile and dump all the
	// concordances

	scanner := bufio.NewScanner(tmpfile)

	for scanner.Scan() {

		var c whosonfirst.WOFConcordances

		raw := scanner.Text()
		dec := json.NewDecoder(strings.NewReader(raw))

		for {

			err := dec.Decode(&c)

			if err == io.EOF {
				break
			}

			if err != nil {
				return err
			}
		}

		out := make([]string, 0)

		for _, k := range header {

			v, ok := c[k]

			if !ok {
				v = ""
			}

			out = append(out, v)
		}

		writer.Write(out)
		writer.Flush()
	}

	writer.Flush()
	return nil
}

func load_feature(fh io.Reader, ctx context.Context) (geojson.Feature, error) {

	ok, err := utils.IsPrincipalWOFRecord(fh, ctx)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	f, err := feature.LoadFeatureFromReader(fh)

	if err != nil {

		path, p_err := index.PathForContext(ctx)

		if p_err != nil {
			msg := fmt.Sprintf("%s (failed to determine path for filehandle because %s)", err, p_err)
			return nil, errors.New(msg)
		}

		msg := fmt.Sprintf("failed to load %s because %s", path, err)
		return nil, errors.New(msg)
	}

	return f, nil
}
