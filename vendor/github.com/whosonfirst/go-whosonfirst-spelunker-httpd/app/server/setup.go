package server

import (
	"context"
	"fmt"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
)

func setupCommon() {

	ctx := context.Background()
	var err error

	sp, err = spelunker.NewSpelunker(ctx, spelunker_uri)

	if err != nil {
		setupCommonError = fmt.Errorf("Failed to set up network, %w", err)
	}
}
