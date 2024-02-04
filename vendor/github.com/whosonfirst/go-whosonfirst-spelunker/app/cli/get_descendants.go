package cli

import (
	"context"
	"fmt"
	"log/slog"

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

	r, _, err := sp.GetDescendants(ctx, id, pg_opts)

	if err != nil {
		return fmt.Errorf("Failed to retrieve descendants, %w", err)
	}

	slog.Info("OK", "count", len(r))
	return nil
}
