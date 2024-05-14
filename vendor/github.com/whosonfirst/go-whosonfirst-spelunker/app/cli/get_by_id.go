package cli

import (
	"context"
	"fmt"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

func get_by_id(ctx context.Context, sp spelunker.Spelunker) error {

	uri_args := new(uri.URIArgs)

	body, err := sp.GetRecordForId(ctx, id, uri_args)

	if err != nil {
		return fmt.Errorf("Failed to get record by ID, %w", err)
	}

	fmt.Println(string(body))
	return nil
}
