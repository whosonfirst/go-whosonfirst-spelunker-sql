package www

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/aaronland/go-http-sanitize"
	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
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
	PageTitle     string
	URIs          *httpd.URIs
	Places        []spr.StandardPlacesResult
	Pagination    pagination.Results
	PaginationURL string
}

func SearchHandler(opts *SearchHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("search")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'search' template")
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
			http.Error(rsp, "womp womp", http.StatusInternalServerError)
			return
		}

		if q != "" {

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

			search_opts := &spelunker.SearchOptions{
				Query: q,
			}

			r, pg_r, err := opts.Spelunker.Search(ctx, pg_opts, search_opts)

			if err != nil {
				logger.Error("Failed to get search", "error", err)
				http.Error(rsp, "womp womp", http.StatusInternalServerError)
				return
			}

			vars.Places = r.Results()
			vars.Pagination = pg_r

			pagination_q := &url.Values{}
			pagination_q.Set("q", q)

			pagination_url := opts.URIs.Search + "?" + pagination_q.Encode()
			vars.PaginationURL = pagination_url
		}

		rsp.Header().Set("Content-Type", "text/html")

		err = t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to return ", "error", err)
			http.Error(rsp, "womp womp", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
