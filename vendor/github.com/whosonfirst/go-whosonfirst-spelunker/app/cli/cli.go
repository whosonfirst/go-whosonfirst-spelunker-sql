package cli

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/sfomuseum/go-flags/flagset"
	spelunker "github.com/whosonfirst/go-whosonfirst-spelunker"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *slog.Logger) error {

	flagset.Parse(fs)

	slog.SetDefault(logger)

	sp, err := spelunker.NewSpelunker(ctx, spelunker_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new spelunker, %w", err)
	}

	switch command {
	case "descendants":
		return get_descendants(ctx, sp)
	case "id":
		return get_by_id(ctx, sp)
	default:
		return fmt.Errorf("Invalid or unsupported command")
	}

	return nil
}
