package cli

import (
	"context"
	"fmt"

	spelunker "github.com/whosonfirst/go-whosonfirst-spelunker"
)

func get_by_id(ctx context.Context, sp spelunker.Spelunker) error {

	body, err := sp.GetById(ctx, id)

	if err != nil {
		return fmt.Errorf("Failed to get record by ID, %w", err)
	}

	fmt.Println(string(body))
	return nil
}
