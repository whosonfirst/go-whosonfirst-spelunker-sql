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

type PlacetypesHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type PlacetypesHandlerVars struct {
	PageTitle string
	URIs      *httpd.URIs
	Facets    []*spelunker.FacetCount
}

func PlacetypesHandler(opts *PlacetypesHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("placetypes")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'placetypes' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("request", req.URL)

		faceting, err := opts.Spelunker.GetPlacetypes(ctx)

		if err != nil {
			logger.Error("Failed to get placetypes", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		vars := PlacetypesHandlerVars{
			PageTitle: "Placetypes",
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
