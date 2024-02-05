package main

import (
	"context"
	"log/slog"
	"flag"
	"os"
	"net/http"
	"net/url"
	"path/filepath"
	"github.com/psanford/sqlite3vfs"
	"github.com/psanford/sqlite3vfshttp"
	
	_ "github.com/mattn/go-sqlite3"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/app/server"
	_ "github.com/whosonfirst/go-whosonfirst-spelunker-sql"
)

type roundTripper struct {
	referer   string
	userAgent string
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if rt.referer != "" {
		req.Header.Set("Referer", rt.referer)
	}

	if rt.userAgent != "" {
		req.Header.Set("User-Agent", rt.userAgent)
	}

	tr := http.DefaultTransport

	if req.URL.Scheme == "file" {
		path := req.URL.Path
		root := filepath.Dir(path)
		base := filepath.Base(path)
		tr = http.NewFileTransport(http.Dir(root))
		req.URL.Path = base
	}

	return tr.RoundTrip(req)
}

func main() {

	ctx := context.Background()
	logger := slog.Default()

	opts := &server.RunOptions{
		Logger: logger,
	}
	
	fs := server.DefaultFlagSet()

	opts, err := server.RunOptionsFromFlagSet(ctx, fs, logger)

	if err != nil {
		logger.Error("Failed to derive run options", "error", err)
		os.Exit(1)
	}
	
	fs.VisitAll(func(fl *flag.Flag){
		
		if fl.Name != "spelunker-uri" {
			return
		}

		spelunker_uri := fl.Value.String()		
		u, err := url.Parse(spelunker_uri)

		if err != nil {
			slog.Error("Failed to parse spelunker URI", "error", err)
			return
		}

		if u.Host != "sqlite3" {
			return
		}

		q := u.Query()

		if !q.Has("vfs") {
			return
		}
		
		vfs_url := q.Get("vfs")
		
		vfs := sqlite3vfshttp.HttpVFS{
			URL:          vfs_url,
			RoundTripper: &roundTripper{
				// referer:   *referer,
				// userAgent: *userAgent,
			},
		}
		
		err = sqlite3vfs.RegisterVFS("httpvfs", &vfs)

		if err != nil {
			slog.Error("Failed to register VFS", "error", err)
			return
		}

		dsn := "not_a_real_name.db?vfs=httpvfs&mode=ro"
		q.Set("dsn", dsn)
		q.Del("vfs")
		
		u.RawQuery = q.Encode()
		opts.SpelunkerURI = u.String()
	})

	err = server.RunWithOptions(ctx, opts)

	if err != nil {
		slog.Error("Failed to run server", "error", err)
		os.Exit(1)
	}
}

