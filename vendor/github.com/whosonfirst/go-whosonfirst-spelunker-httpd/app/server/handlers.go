package server

import (
	"context"
	"fmt"
	"net/http"
	"log/slog"
	
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/www"
)

func geoJSONHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	opts := &www.GeoJSONHandlerOptions{
		Spelunker: sp,
	}

	return www.GeoJSONHandler(opts)
}
