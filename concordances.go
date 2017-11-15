package concordances

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-index/utils"
	"io"
	_ "log"
	"sync"
)

func load_feature(fh io.Reader, ctx context.Context) (geojson.Feature, error) {

	ok, err := utils.IsPrincipalWOFRecord(fh, ctx)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	f, err := feature.LoadFeatureFromReader(fh)

	// log.Println(f)

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

	possible, err := ListConcordances(mode, sources...)

	if err != nil {
		return err
	}

	writer := csv.NewWriter(out)
	writer.Write(possible)
	writer.Flush()

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

		row := make([]string, 0)
		matches := 0

		for _, key := range possible {

			value, ok := c[key]

			if ok {
				matches += 1
			} else {
				value = ""
			}

			row = append(row, value)
		}

		if matches > 1 { // wof:id
			mu.Lock()
			writer.Write(row)
			mu.Unlock()
		}

		return nil
	}

	idx, err := index.NewIndexer(mode, cb)

	if err != nil {
		return err
	}

	for _, src := range sources {

		err := idx.IndexPath(src)

		if err != nil {
			return err
		}
	}

	return nil
}
