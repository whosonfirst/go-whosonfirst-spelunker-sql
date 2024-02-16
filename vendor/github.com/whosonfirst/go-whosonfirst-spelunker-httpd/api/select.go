package api

import (
	"log/slog"
	"net/http"
	"regexp"

	"encoding/json"
	"github.com/aaronland/go-http-sanitize"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type SelectHandlerOptions struct {
	Pattern   *regexp.Regexp
	Spelunker spelunker.Spelunker
}

func SelectHandler(opts *SelectHandlerOptions) (http.Handler, error) {

	logger := slog.Default()

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger = logger.With("request", req.URL)
		logger = logger.With("address", req.RemoteAddr)

		query, err := sanitize.GetString(req, "select")

		if err != nil {
			http.Error(rsp, "Invalid parameters", http.StatusBadRequest)
			return
		}

		if query == "" {
			http.Error(rsp, "Missing select", http.StatusBadRequest)
			return
		}

		if !opts.Pattern.MatchString(query) {
			http.Error(rsp, "Invalid select", http.StatusBadRequest)
			return
		}

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

		query_rsp := gjson.GetBytes(r, query)

		var rsp_body []byte

		if query_rsp.Exists() {

			enc, err := json.Marshal(query_rsp.Value())

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			rsp_body = enc
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Write(rsp_body)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
