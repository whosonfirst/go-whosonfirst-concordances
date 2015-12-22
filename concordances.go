package concordances

import (
	crawl "github.com/whosonfirst/go-whosonfirst-crawl"
	geojson "github.com/whosonfirst/go-whosonfirst-geojson"
	"os"
	"sync"
)

func ListConcordances(root string) []string {

	tmp := make(map[string]int)
	wg := new(sync.WaitGroup)

	callback := func(source string, info os.FileInfo) error {

		wg.Add(1)

		defer wg.Done()

		if info.IsDir() {
			return nil
		}

		feature, err := geojson.UnmarshalFile(source)

		if err != nil {
			return err
		}

		body := feature.Body()
		props, _ := body.S("properties").ChildrenMap()

		for key, child := range props {

			if key != "wof:concordances" {
				continue
			}

			concordances, _ := child.ChildrenMap()

			for src, _ := range concordances {
				tmp[src] += 1
			}

			break
		}

		return nil
	}

	c := crawl.NewCrawler(root)
	_ = c.Crawl(callback)

	wg.Wait()

	concordances := make([]string, 0)

	for name, _ := range tmp {
		concordances = append(concordances, name)
	}

	return concordances
}
