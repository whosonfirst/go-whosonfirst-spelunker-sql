package www

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sfomuseum/go-http-auth"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type IdHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type IdHandlerAncestor struct {
	Placetype string
	Id        int64
}

type IdHandlerVars struct {
	Id               int64
	RequestId        string
	URIArgs          *uri.URIArgs
	PageTitle        string
	URIs             *httpd.URIs
	Properties       string
	CountDescendants int64
	Hierarchies      [][]*IdHandlerAncestor
	RelPath          string
	GitHubURL        string
	WriteFieldURL    string
}

func IdHandler(opts *IdHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("id")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'id' template")
	}

	alt_t := opts.Templates.Lookup("alt")

	if alt_t == nil {
		return nil, fmt.Errorf("Missing alt template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("request", req.URL)

		req_uri, err, status := httpd.ParseURIFromRequest(req, nil)

		if err != nil {
			slog.Error("Failed to parse URI from request", "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), status)
			return
		}

		wof_id := req_uri.Id

		req_id, err := uri.Id2Fname(req_uri.Id, req_uri.URIArgs)

		if err != nil {
			slog.Error("Failed to derive request ID", "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		req_id = strings.Replace(req_id, filepath.Ext(req_id), "", 1)

		logger = logger.With("request id", req_id)
		logger = logger.With("wof id", wof_id)

		f, err := httpd.FeatureFromRequestURI(ctx, opts.Spelunker, req_uri)

		if err != nil {
			slog.Error("Failed to get by ID", "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		props := gjson.GetBytes(f, "properties")
		page_title := gjson.GetBytes(f, "properties.wof:name")

		rel_path, err := uri.Id2RelPath(wof_id, req_uri.URIArgs)

		if err != nil {
			slog.Error("Failed to derive relative path for record", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		repo_name := gjson.GetBytes(f, "properties.wof:repo")

		github_url := fmt.Sprintf("https://github.com/whosonfirst-data/%s/blob/master/data/%s", repo_name, rel_path)

		vars := IdHandlerVars{
			Id:         wof_id,
			RequestId:  req_id,
			URIArgs:    req_uri.URIArgs,
			Properties: props.String(),
			PageTitle:  page_title.String(),
			GitHubURL:  github_url,
			URIs:       opts.URIs,
			RelPath:    rel_path,
		}

		if req_uri.IsAlternate {

			rsp.Header().Set("Content-Type", "text/html")

			err = alt_t.Execute(rsp, vars)

			if err != nil {
				slog.Error("Failed to return ", "error", err)
				http.Error(rsp, "womp womp", http.StatusInternalServerError)
			}

			return
		}

		count_descendants, err := opts.Spelunker.CountDescendants(ctx, wof_id)

		if err != nil {
			slog.Error("Failed to count descendants", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		// START OF there's got to be a better way to do this...

		str_pt := gjson.GetBytes(f, "properties.wof:placetype")

		pt, err := placetypes.GetPlacetypeByName(str_pt.String())

		if err != nil {
			slog.Error("Failed to load placetype", "placetype", str_pt, "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		roles := []string{
			"common",
			"optional",
			"common_optional",
		}

		ancestors := placetypes.AncestorsForRoles(pt, roles)
		count_ancestors := len(ancestors)

		sorted := make([]string, 0)

		for i := count_ancestors - 1; i >= 0; i-- {
			n := ancestors[i]
			sorted = append(sorted, n.String())
		}

		hierarchies := properties.Hierarchies(f)

		handler_hierarchies := make([][]*IdHandlerAncestor, len(hierarchies))

		for idx, hier := range hierarchies {

			handler_ancestors := make([]*IdHandlerAncestor, 0)

			for _, n := range sorted {

				k := fmt.Sprintf("%s_id", n)
				v, ok := hier[k]

				if !ok {
					continue
				}

				a := &IdHandlerAncestor{
					Placetype: n,
					Id:        v,
				}

				handler_ancestors = append(handler_ancestors, a)
			}

			handler_hierarchies[idx] = handler_ancestors
		}

		// END OF there's got to be a better way to do this...

		writefield_url := fmt.Sprintf("https://raw.githubusercontent.com/whosonfirst-data/%s/master/data/%s", repo_name, rel_path)

		vars.CountDescendants = count_descendants
		vars.Hierarchies = handler_hierarchies
		vars.WriteFieldURL = writefield_url

		rsp.Header().Set("Content-Type", "text/html")

		err = t.Execute(rsp, vars)

		if err != nil {
			slog.Error("Failed to return ", "error", err)
			http.Error(rsp, "womp womp", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
