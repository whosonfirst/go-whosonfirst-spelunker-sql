package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type SPRHandlerOptions struct {
	Spelunker spelunker.Spelunker
}

func SPRHandler(opts *SPRHandlerOptions) (http.Handler, error) {

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

		s, err := spr.WhosOnFirstSPR(r)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(s)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
