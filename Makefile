CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-concordances; then rm -rf src/github.com/whosonfirst/go-whosonfirst-concordances; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-concordances
	cp concordances.go src/github.com/whosonfirst/go-whosonfirst-concordances/concordances.go
	cp -r vendor/src/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-crawl"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor/src; then rm -rf vendor/src; fi
	cp -r src vendor/src
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go
	go fmt cmd/*.go


bin:	self
	@GOPATH=$(GOPATH) go build -o bin/wof-concordances-list cmd/wof-concordances-list.go
	@GOPATH=$(GOPATH) go build -o bin/wof-concordances-write cmd/wof-concordances-write.go
