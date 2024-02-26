package www

import (
	"encoding/json"
	"log/slog"
	"net/http"

	// "github.com/aaronland/go-pagination"
	// "github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type DescendantsFacetHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
}

func DescendantsFacetHandler(opts *DescendantsFacetHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("request", req.URL)

		uri, err, status := httpd.ParseURIFromRequest(req, nil)

		if err != nil {
			logger.Error("Failed to parse URI from request", "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), status)
			return
		}

		logger = logger.With("wofid", uri.Id)

		/*
			pg_opts, err := countable.NewCountableOptions()

			if err != nil {
				logger.Error("Failed to create pagination options", "error", err)
				http.Error(rsp, "womp womp", http.StatusInternalServerError)
				return
			}

			pg, pg_err := httpd.ParsePageNumberFromRequest(req)

			if pg_err == nil {
				pg_opts.Pointer(pg)
			}
		*/

		filter_params := []string{
			"placetype",
			"country",
		}

		filters, err := FiltersFromRequest(ctx, req, filter_params)

		if err != nil {
			logger.Error("Failed to derive filters from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		facets, err := FacetsFromRequest(ctx, req, filter_params)

		if err != nil {
			logger.Error("Failed to derive facets from requrst", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if len(facets) == 0 {
			logger.Error("No facets from requrst")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		facets_rsp, err := opts.Spelunker.FacetDescendants(ctx, uri.Id, filters, facets)

		if err != nil {
			logger.Error("Failed to get facets for descendants", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(facets_rsp)

		if err != nil {
			logger.Error("Failed to encode facets response", "error", err)
			http.Error(rsp, "womp womp", http.StatusInternalServerError)
			return
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
