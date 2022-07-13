package concordances

import (
	"context"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"io"
	"sort"
	"sync"
)

// ListKeys() returns the list of unique keys for all the concordances found in 'iterator_sources'.
// 'iterator_uri' is expected to be a valid `whosonfirst/go-whosonfirst-iterate/v2` URI and 'iterator_sources'
// a list of URIs to be crawled.
func ListKeys(ctx context.Context, iterator_uri string, iterator_sources ...string) ([]string, error) {

	sources := new(sync.Map)

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", path, err)
		}

		c := properties.Concordances(body)

		for src, _ := range c {
			sources.Store(src, true)
		}

		return nil
	}

	iter, err := iterator.NewIterator(ctx, iterator_uri, iter_cb)

	if err != nil {
		return nil, fmt.Errorf("Failed to create iterator, %w", err)
	}

	err = iter.IterateURIs(ctx, iterator_sources...)

	if err != nil {
		return nil, fmt.Errorf("Failed to err iterate URIs, %w", err)
	}

	concordances_keys := make([]string, 0)

	sources.Range(func(k interface{}, v interface{}) bool {
		concordances_keys = append(concordances_keys, k.(string))
		return true
	})

	sort.Strings(concordances_keys)
	return concordances_keys, nil
}
