prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-concordances; then rm -rf src/github.com/whosonfirst/go-whosonfirst-concordances; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-concordances
	cp concordances.go src/github.com/whosonfirst/go-whosonfirst-concordances/concordances.go

deps:
	@GOPATH=$(shell pwd) \
	go get -u "github.com/whosonfirst/go-whosonfirst-crawl"
	@GOPATH=$(shell pwd) \
	go get -u "github.com/whosonfirst/go-whosonfirst-geojson"

fmt:
	go fmt *.go
	go fmt cmd/*.go


bin:	self
	@GOPATH=$(shell pwd) \
	go build -o bin/wof-concordances-list cmd/wof-concordances-list.go
	@GOPATH=$(shell pwd) \
	go build -o bin/wof-concordances-write cmd/wof-concordances-write.go
