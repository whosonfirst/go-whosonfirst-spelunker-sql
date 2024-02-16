package cli

import (
	"context"
	"encoding/json"
	"fmt"
	_ "log/slog"
	"os"

	"github.com/aaronland/go-pagination/countable"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
)

func search(ctx context.Context, sp spelunker.Spelunker) error {

	// Eventually we'll need to check if we're doing cursor-base pagination

	pg_opts, err := countable.NewCountableOptions()

	if err != nil {
		return fmt.Errorf("Failed to create countable options, %w", err)
	}

	pg_opts.PerPage(per_page)
	pg_opts.Pointer(page)

	search_opts := &spelunker.SearchOptions{
		Query: query,
	}

	r, _, err := sp.Search(ctx, pg_opts, search_opts)

	if err != nil {
		return fmt.Errorf("Failed to retrieve descendants, %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	err = enc.Encode(r)

	if err != nil {
		return fmt.Errorf("Failed to encode descendants, %w", err)
	}

	return nil
}
