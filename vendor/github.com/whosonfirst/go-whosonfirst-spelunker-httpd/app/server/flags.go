package server

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var server_uri string
var spelunker_uri string
var authenticator_uri string
var protomaps_api_key string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("spelunker")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "...")
	fs.StringVar(&spelunker_uri, "spelunker-uri", "null://", "...")
	fs.StringVar(&authenticator_uri, "authenticator-uri", "null://", "...")
	fs.StringVar(&protomaps_api_key, "protomaps-api-key", "", "...")
	return fs
}
