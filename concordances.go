package concordances

import (
	"encoding/csv"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-crawl"
	"io"
	"io/ioutil"
	_ "log"
	"os"
	"sync"
)

type CrawlFunc func(concordance map[string]string) error

func ListConcordances(root string) ([]string, error) {

	tmp := make(map[string]int)
	concordances := make([]string, 0)

	mu := sync.Mutex{}

	dothis := func(concordances map[string]string) error {

		mu.Lock()

		for src, _ := range concordances {
			tmp[src] += 1
		}

		mu.Unlock()
		return nil
	}

	err := CrawlConcordances(root, dothis)

	if err != nil {
		return concordances, err
	}

	for name, _ := range tmp {
		concordances = append(concordances, name)
	}

	return concordances, nil
}

func WriteConcordances(root string, out io.Writer) error {

	possible, err := ListConcordances(root)

	if err != nil {
		return err
	}

	writer := csv.NewWriter(out)
	writer.Write(possible)
	writer.Flush()

	mu := sync.Mutex{}

	dothis := func(concordances map[string]string) error {

		row := make([]string, 0)
		matches := 0

		for _, key := range possible {

			value, ok := concordances[key]

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

	err = CrawlConcordances(root, dothis)

	if err != nil {
		return err
	}

	writer.Flush()
	return nil
}

func CrawlConcordances(root string, dothis CrawlFunc) error {

	callback := func(source string, info os.FileInfo) error {

		if info.IsDir() {
			return nil
		}

		concordances, err := LoadConcordances(source)

		if err != nil {
			return err
		}

		return dothis(concordances)
	}

	c := crawl.NewCrawler(root)
	return c.Crawl(callback)
}

func LoadConcordances(path string) (map[string]string, error) {

	concordances := make(map[string]string)

	fh, err := os.Open(path)

	if err != nil {
		return concordances, err
	}

	feature, err := ioutil.ReadAll(fh)

	if err != nil {
		return concordances, err
	}

	r := gjson.GetBytes(feature, "properties.wof:concordances")

	if !r.Exists() {
		return concordances, errors.New("Feature missing a wof:concordances property")
	}

	for k, v := range r.Map() {
		concordances[k] = v.String()
	}

	r = gjson.GetBytes(feature, "properties.wof:id")
	concordances["wof:id"] = r.String()

	return concordances, nil
}
