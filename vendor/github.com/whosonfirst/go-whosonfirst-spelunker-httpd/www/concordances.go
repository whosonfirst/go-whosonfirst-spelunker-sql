package www

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type ConcordancesHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type ConcordancesHandlerVars struct {
	PageTitle string
	URIs      *httpd.URIs
	Facets    []*spelunker.FacetCount
}

func ConcordancesHandler(opts *ConcordancesHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("concordances")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'concordances' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("request", req.URL)

		faceting, err := opts.Spelunker.GetConcordances(ctx)

		if err != nil {
			logger.Error("Failed to get concordances", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		vars := ConcordancesHandlerVars{
			PageTitle: "Concordances",
			URIs:      opts.URIs,
			Facets:    faceting.Results,
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
