package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/whosonfirst/go-whosonfirst-spelunker-sql"
	"github.com/whosonfirst/go-whosonfirst-spelunker/app/cli"
)

func main() {

	ctx := context.Background()
	err := cli.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
