CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -d src/github.com/whosonfirst/go-whosonfirst-concordances; then rm -rf src/github.com/whosonfirst/go-whosonfirst-concordances; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-concordances
	cp concordances.go src/github.com/whosonfirst/go-whosonfirst-concordances/concordances.go
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi

deps:
	@GOPATH=$(GOPATH) go get -u "github.com/facebookgo/atomicfile"
	@GOPATH=$(GOPATH) go get -u "github.com/tidwall/gjson"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-crawl"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-index"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-repo"

vendor-deps: rmdeps deps
	if test -d vendor; then rm -rf vendor; fi
	mkdir vendor
	cp -r src/* vendor/
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt *.go
	go fmt cmd/*.go


bin:	self
	@GOPATH=$(GOPATH) go build -o bin/wof-concordances-list cmd/wof-concordances-list.go
	@GOPATH=$(GOPATH) go build -o bin/wof-concordances-write cmd/wof-concordances-write.go
	@GOPATH=$(GOPATH) go build -o bin/wof-build-concordances cmd/wof-build-concordances.go

dist: self
	OS=darwin make dist-os
	OS=windows make dist-os
	OS=linux make dist-os

dist-os:
	mkdir -p dist/$(OS)
	GOOS=$(OS) GOPATH=$(GOPATH) GOARCH=386 go build -o dist/$(OS)/wof-build-concordances cmd/wof-build-concordances.go
	cd dist/$(OS) && shasum -a 256 wof-build-concordances > wof-build-concordances.sha256

rmdist:
	if test -d dist; then rm -rf dist; fi
