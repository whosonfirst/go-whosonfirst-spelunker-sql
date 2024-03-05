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

func get_descendants(ctx context.Context, sp spelunker.Spelunker) error {

	// Eventually we'll need to check if we're doing cursor-base pagination

	pg_opts, err := countable.NewCountableOptions()

	if err != nil {
		return fmt.Errorf("Failed to create countable options, %w", err)
	}

	pg_opts.PerPage(per_page)
	pg_opts.Pointer(page)

	filters := make([]spelunker.Filter, 0)

	r, _, err := sp.GetDescendants(ctx, pg_opts, id, filters)

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
