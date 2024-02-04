package cli

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

var spelunker_uri string
var command string

var id int64

var per_page int64

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("descendants")
	fs.StringVar(&spelunker_uri, "spelunker-uri", "", "...")

	fs.StringVar(&command, "command", "", "...")
	fs.Int64Var(&id, "id", 0, "...")

	fs.Int64Var(&per_page, "per-page", 10, "...")	

	return fs
}
