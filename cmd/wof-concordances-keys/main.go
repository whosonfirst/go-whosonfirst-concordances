package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-concordances"
	_ "github.com/whosonfirst/go-whosonfirst-iterate-organization"
	"log"
	"os"
	"strings"
)

func main() {

	iterator_uri := flag.String("iterator-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v2 URI.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "wof-concordances-keys returns the list of unique keys for all the concordances found in one or more sources.")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s source(N) source(N)\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	sources := flag.Args()

	ctx := context.Background()

	list, err := concordances.ListKeys(ctx, *iterator_uri, sources...)

	if err != nil {
		log.Fatalf("Failed to list concordances, %v", err)
	}

	fmt.Println(strings.Join(list, "\n"))
}
