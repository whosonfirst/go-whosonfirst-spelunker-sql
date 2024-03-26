package www

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	"github.com/aaronland/go-pagination"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type SearchHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type SearchHandlerVars struct {
	PageTitle        string
	URIs             *httpd.URIs
	Places           []spr.StandardPlacesResult
	Pagination       pagination.Results
	PaginationURL    string
	FacetsURL        string
	FacetsContextURL string
	SearchOptions    *spelunker.SearchOptions
}

func SearchHandler(opts *SearchHandlerOptions) (http.Handler, error) {

	form_t := opts.Templates.Lookup("search")

	if form_t == nil {
		return nil, fmt.Errorf("Failed to locate 'search' template")
	}

	results_t := opts.Templates.Lookup("search_results")

	if results_t == nil {
		return nil, fmt.Errorf("Failed to locate 'search_results' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("request", req.URL)

		vars := SearchHandlerVars{
			URIs:      opts.URIs,
			PageTitle: "Search",
		}

		q, err := sanitize.GetString(req, "q")

		if err != nil {
			logger.Error("Failed to determine query string", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		if q == "" {

			rsp.Header().Set("Content-Type", "text/html")

			err = form_t.Execute(rsp, vars)

			if err != nil {
				logger.Error("Failed to return ", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			}

			return
		}

		pg_opts, err := httpd.PaginationOptionsFromRequest(req)

		if err != nil {
			logger.Error("Failed to create pagination options", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		search_opts := &spelunker.SearchOptions{
			Query: q,
		}

		filter_params := httpd.DefaultFilterParams()

		filters, err := httpd.FiltersFromRequest(ctx, req, filter_params)

		if err != nil {
			logger.Error("Failed to derive filters from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		r, pg_r, err := opts.Spelunker.Search(ctx, pg_opts, search_opts, filters)

		if err != nil {
			logger.Error("Failed to get search", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		vars.Places = r.Results()
		vars.Pagination = pg_r

		pagination_url := httpd.URIForSearch(opts.URIs.Search, q, filters, nil)
		facets_url := httpd.URIForSearch(opts.URIs.SearchFaceted, q, filters, nil)
		facets_context_url := httpd.URIForSearch(opts.URIs.Search, q, filters, nil)

		vars.PaginationURL = pagination_url
		vars.FacetsURL = facets_url
		vars.FacetsContextURL = facets_context_url
		vars.SearchOptions = search_opts

		rsp.Header().Set("Content-Type", "text/html")

		err = results_t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to return ", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
