package www

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type HasPlacetypeHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type HasPlacetypeHandlerVars struct {
	PageTitle        string
	URIs             *httpd.URIs
	Placetype        *placetypes.WOFPlacetype
	Places           []spr.StandardPlacesResult
	Pagination       pagination.Results
	PaginationURL    string
	FacetsURL        string
	FacetsContextURL string
}

func HasPlacetypeHandler(opts *HasPlacetypeHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("placetype")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'placetype' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("request", req.URL)

		req_pt := req.PathValue("placetype")

		logger = logger.With("request placetype", req_pt)

		pt, err := placetypes.GetPlacetypeByName(req_pt)

		if err != nil {
			logger.Error("Invalid placetype", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		pg_opts, err := countable.NewCountableOptions()

		if err != nil {
			logger.Error("Failed to create pagination options", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		pg, pg_err := httpd.ParsePageNumberFromRequest(req)

		if pg_err == nil {
			logger = logger.With("page", pg)
			pg_opts.Pointer(pg)
		}

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

		r, pg_r, err := opts.Spelunker.HasPlacetype(ctx, pg_opts, pt, filters)

		if err != nil {
			logger.Error("Failed to get records having placetype", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		pagination_url := fmt.Sprintf("%s?", req.URL.Path)

		vars := HasPlacetypeHandlerVars{
			PageTitle:     pt.Name,
			URIs:          opts.URIs,
			Placetype:     pt,
			Places:        r.Results(),
			Pagination:    pg_r,
			PaginationURL: pagination_url,
		}

		rsp.Header().Set("Content-Type", "text/html")

		err = t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to render template", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
