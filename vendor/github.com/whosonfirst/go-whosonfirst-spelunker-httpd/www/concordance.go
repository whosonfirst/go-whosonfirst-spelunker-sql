package www

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-sources"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type HasConcordanceHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type HasConcordanceHandlerVars struct {
	PageTitle        string
	URIs             *httpd.URIs
	Concordance      *spelunker.Concordance
	Places           []spr.StandardPlacesResult
	Pagination       pagination.Results
	PaginationURL    string
	FacetsURL        string
	FacetsContextURL string
	Source           *sources.WOFSource
}

func HasConcordanceHandler(opts *HasConcordanceHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("concordance")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'concordance' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("request", req.URL)

		ns := req.PathValue("namespace")
		pred := req.PathValue("predicate")
		value := req.PathValue("value")

		ns = strings.TrimRight(ns, ":")
		pred = strings.TrimLeft(pred, ":")
		pred = strings.TrimRight(pred, "=")

		c := spelunker.NewConcordanceFromTriple(ns, pred, value)

		logger = logger.With("namespace", ns)
		logger = logger.With("predicate", pred)
		logger = logger.With("value", value)

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

		filters, err := httpd.FiltersFromRequest(ctx, req, filter_params)

		if err != nil {
			logger.Error("Failed to derive filters from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		r, pg_r, err := opts.Spelunker.HasConcordance(ctx, pg_opts, ns, pred, value, filters)

		if err != nil {
			logger.Error("Failed to get records having concordance", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		pagination_url := fmt.Sprintf("%s?", req.URL.Path)

		page_title := fmt.Sprintf("Concordances for %s", c)

		src, err := sources.GetSourceByName(ns)

		if err != nil {
			logger.Error("Failed to derive source from namespace", "error", err)
		}

		vars := HasConcordanceHandlerVars{
			PageTitle:     page_title,
			URIs:          opts.URIs,
			Concordance:   c,
			Places:        r.Results(),
			Pagination:    pg_r,
			PaginationURL: pagination_url,
			Source:        src,
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
