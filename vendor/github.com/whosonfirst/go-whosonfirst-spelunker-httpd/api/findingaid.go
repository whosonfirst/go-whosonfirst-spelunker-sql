package api

import (
	"log/slog"
	"net/http"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type FindingAidHandlerOptions struct {
	Spelunker spelunker.Spelunker
}

func FindingAidHandler(opts *FindingAidHandlerOptions) (http.Handler, error) {

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

		spr, err := httpd.SPRFromRequestURI(ctx, opts.Spelunker, req_uri)

		if err != nil {
			slog.Error("Failed to get by ID", "id", req_uri.Id, "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		repo := spr.Repo()

		rsp.Header().Set("Content-Type", "text/plain")
		rsp.Write([]byte(repo))
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
