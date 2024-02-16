package api

import (
	"log/slog"
	"net/http"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type GeoJSONHandlerOptions struct {
	Spelunker spelunker.Spelunker
}

func GeoJSONHandler(opts *GeoJSONHandlerOptions) (http.Handler, error) {

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

		wof_id := req_uri.Id
		logger = logger.With("wof id", wof_id)

		r, err := httpd.FeatureFromRequestURI(ctx, opts.Spelunker, req_uri)

		if err != nil {
			slog.Error("Failed to get by ID", "id", wof_id, "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Write(r)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
