package api

import (
	"log/slog"
	"net/http"

	"github.com/sfomuseum/go-geojsonld"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type GeoJSONLDHandlerOptions struct {
	Spelunker spelunker.Spelunker
}

func GeoJSONLDHandler(opts *GeoJSONLDHandlerOptions) (http.Handler, error) {

	logger := slog.Default()

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger = logger.With("request", req.URL)
		logger = logger.With("address", req.RemoteAddr)

		req_uri, err, status := httpd.ParseURIFromRequest(req, nil)

		if err != nil {
			slog.Error("Failed to parse URI from request", "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), status)
			return
		}

		r, err := httpd.FeatureFromRequestURI(ctx, opts.Spelunker, req_uri)

		if err != nil {
			slog.Error("Failed to get by ID", "id", req_uri.Id, "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		body, err := geojsonld.AsGeoJSONLD(ctx, r)

		if err != nil {
			slog.Error("Failed to render geojson", "id", req_uri.Id, "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/geo+json")
		rsp.Write([]byte(body))
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
