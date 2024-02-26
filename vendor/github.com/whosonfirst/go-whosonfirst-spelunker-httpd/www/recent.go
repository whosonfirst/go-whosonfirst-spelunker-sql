package www

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/sfomuseum/go-http-auth"
	"github.com/sfomuseum/iso8601duration"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type RecentHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type RecentHandlerVars struct {
	PageTitle     string
	URIs          *httpd.URIs
	Places        []spr.StandardPlacesResult
	Pagination    pagination.Results
	PaginationURL string
	Duration      time.Duration
}

func RecentHandler(opts *RecentHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("recent")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'recent' template")
	}

	re_full, err := regexp.Compile(`P((?P<year>\d+)Y)?((?P<month>\d+)M)?((?P<day>\d+)D)?(T((?P<hour>\d+)H)?((?P<minute>\d+)M)?((?P<second>\d+)S)?)?`)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile ISO8601 duration pattern, %w", err)
	}

	re_week, err := regexp.Compile(`P((?P<week>\d+)W)`)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile ISO8601 duration week pattern, %w", err)
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("request", req.URL)

		slog.Info("Get recent")

		str_d := "P30D"

		path := req.URL.Path

		if re_week.MatchString(path) {
			m := re_week.FindStringSubmatch(path)
			str_d = m[0]
		} else if re_full.MatchString(path) {
			m := re_full.FindStringSubmatch(path)
			str_d = m[0]
		} else {
			// pass
		}

		logger = logger.With("duration", str_d)

		d, err := duration.FromString(str_d)

		if err != nil {
			logger.Error("Failed to parse duration", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

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

		filters := make([]spelunker.Filter, 0)

		r, pg_r, err := opts.Spelunker.GetRecent(ctx, pg_opts, d.ToDuration(), filters)

		if err != nil {
			logger.Error("Failed to get recent", "error", err)
			http.Error(rsp, "womp womp", http.StatusInternalServerError)
			return
		}

		pagination_url := req.URL.Path

		vars := RecentHandlerVars{
			Places:        r.Results(),
			Pagination:    pg_r,
			URIs:          opts.URIs,
			PaginationURL: pagination_url,
			Duration:      d.ToDuration(),
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
