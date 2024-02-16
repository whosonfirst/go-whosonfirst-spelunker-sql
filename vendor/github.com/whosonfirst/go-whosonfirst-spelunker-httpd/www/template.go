package www

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type TemplateHandlerOptions struct {
	Authenticator auth.Authenticator
	Templates     *template.Template
	TemplateName  string
	PageTitle     string
	URIs          *httpd.URIs
}

type TemplateHandlerVars struct {
	Id         int64
	PageTitle  string
	URIs       *httpd.URIs
	Properties string
}

func TemplateHandler(opts *TemplateHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup(opts.TemplateName)

	if t == nil {
		return nil, fmt.Errorf("Failed to locate ihelp' template")
	}

	logger := slog.Default()

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger = logger.With("request", req.URL)
		logger = logger.With("address", req.RemoteAddr)

		vars := TemplateHandlerVars{
			PageTitle: opts.PageTitle,
			URIs:      opts.URIs,
		}

		rsp.Header().Set("Content-Type", "text/html")

		err := t.Execute(rsp, vars)

		if err != nil {
			slog.Error("Failed to return ", "error", err)
			http.Error(rsp, "womp womp", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
