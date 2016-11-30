package concordances

import (
	"encoding/csv"
	_ "fmt"
	crawl "github.com/whosonfirst/go-whosonfirst-crawl"
	geojson "github.com/whosonfirst/go-whosonfirst-geojson"
	"io"
	_ "log"
	"os"
	"strconv"
	"sync"
	_ "time"
)

type CrawlFunc func(concordance map[string]string)

func ListConcordances(root string) []string {

	tmp := make(map[string]int)
	concordances := make([]string, 0)

	mu := sync.Mutex{}

	dothis := func(concordances map[string]string) {

		mu.Lock()

		for src, _ := range concordances {

			tmp[src] += 1
		}

		mu.Unlock()
	}

	CrawlConcordances(root, dothis)

	for name, _ := range tmp {
		concordances = append(concordances, name)
	}

	return concordances
}

func WriteConcordances(root string, out io.Writer) {

	possible := ListConcordances(root)

	writer := csv.NewWriter(out)
	writer.Write(possible)
	writer.Flush()

	mu := sync.Mutex{}

	dothis := func(concordances map[string]string) {

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
	}

	CrawlConcordances(root, dothis)
	writer.Flush()
}

func CrawlConcordances(root string, dothis CrawlFunc) {

	// wg := new(sync.WaitGroup)

	callback := func(source string, info os.FileInfo) error {

		//  wg.Add(1)
		// defer wg.Done()

		if info.IsDir() {
			return nil
		}

		concordances, err := LoadConcordances(source)

		if err == nil {
			dothis(concordances)
		} else {
			// log.Println(err)
		}

		return nil
	}

	c := crawl.NewCrawler(root)
	_ = c.Crawl(callback)

	// wg.Wait()
}

// please to be caching me... (20151221/thisisaaronland)
// to investigate: https://github.com/patrickmn/go-cache

func LoadConcordances(path string) (map[string]string, error) {

	concordances := make(map[string]string)

	feature, err := geojson.UnmarshalFile(path)

	if err != nil {
		// fmt.Println(source, err)
		return concordances, err
	}

	body := feature.Body()
	props, _ := body.S("properties").ChildrenMap()

	wof_id := feature.Id()
	str_id := strconv.Itoa(wof_id)

	concordances["wof:id"] = str_id

	for key, child := range props {

		if key != "wof:concordances" {
			continue
		}

		possible, _ := child.ChildrenMap()

		for src, id := range possible {

			var str_id string
			var float_id float64
			var ok bool

			str_id, ok = id.Data().(string)

			if ok {
				concordances[src] = str_id
				continue
			}

			float_id, ok = id.Data().(float64)

			if ok {
				str_id := strconv.FormatFloat(float_id, 'f', -1, 64)
				concordances[src] = str_id
				continue
			}

			// fmt.Printf("failed to handle %s=%v\n", src, id)
		}

		break
	}

	return concordances, nil
}
