# go-whosonfirst-concordances

A Go package for working with Who's On First concordances

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-whosonfirst-concordances.svg)](https://pkg.go.dev/github.com/whosonfirst/go-whosonfirst-concordances)

## Tools

```
$> make cli
go build -mod vendor -o bin/wof-concordances-keys cmd/wof-concordances-keys/main.go
```

### wof-concordances-keys

`wof-concordances-keys` returns the list of unique keys for all the concordances found in one or more sources.

```
$> ./bin/wof-concordances-keys -h
wof-concordances-keys returns the list of unique keys for all the concordances found in one or more sources.Usage:
	 ./bin/wof-concordances-keys source(N) source(N)
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v2 URI. (default "repo://")
```

For example:

```
$> ./bin/wof-concordances-keys \
	/usr/local/data/sfomuseum-data-whosonfirst/ \
	/usr/local/data/sfomuseum-data-enterprise/
4sq:id
chsdm:person
dbp:id
digitalenvoy:country_code
digitalenvoy:metro_code
digitalenvoy:region_code
faa:code
fb:id
fct:id
fifa:id
fips:code
flysfo:code
gaul:id
gn:id
gp:id
hasc:id
iata
iata:callsign
iata:code
icao
icao:callsign
icao:code
ioc:id
iso:id
itu:id
loc:id
m49:code
marc:id
mzb:id
ne:adm0_a3
ne:id
nyt:id
oa:id
pl-gugik
qs:id
qs_pg:id
ro-ancpi:id
sg:id
tgn:id
uncrt:id
unlc:id
uscensus:geoid
wd:id
wikidata
wikipedia
wk:id
wk:page
wk:pageid
wmo:id
```

This tools support the [go-whosonfirst-iterate-organization](https://github.com/whosonfirst/go-whosonfirst-iterate-organization) package so you can iterate over all the Who's On First repositories in a given organization. For example, here is how you might list the concordances keys for all the `whosonfirst-data-admin-*` repositories in the [whosonfirst-data](https://github.com/whosonfirst-data/whosonfirst-data/) organization:

```
$> ./bin/wof-concordances-keys \
   -iterator-uri org://tmp \
   'whosonfirst-data://?prefix=whosonfirst-data-admin-'
```

_Note: This example will take a long time to complete._