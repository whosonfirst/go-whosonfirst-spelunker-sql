package main

import (
	"context"
	"fmt"
	"github.com/psanford/sqlite3vfs"
	"github.com/psanford/sqlite3vfshttp"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/app/server"
	_ "github.com/whosonfirst/go-whosonfirst-spelunker-sql"
)

// https://github.com/psanford/sqlite3vfshttp/blob/main/sqlitehttpcli/sqlitehttpcli.go

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

	fs := server.DefaultFlagSet()

	opts, err := server.RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		logger.Error("Failed to derive run options", "error", err)
		os.Exit(1)
	}

	is_vfs, vfs_uri, err := checkVFS(opts.SpelunkerURI)

	if err != nil {
		logger.Error("Failed to parse spelunker URI", "error", err)
		os.Exit(1)
	}

	if is_vfs {
		opts.SpelunkerURI = vfs_uri
	}

	err = server.RunWithOptions(ctx, opts, logger)

	if err != nil {
		slog.Error("Failed to run server", "error", err)
		os.Exit(1)
	}
}

func checkVFS(spelunker_uri string) (bool, string, error) {

	u, err := url.Parse(spelunker_uri)

	if err != nil {
		return false, spelunker_uri, err
	}

	if u.Host != "sqlite3" {
		return false, spelunker_uri, nil
	}

	q := u.Query()

	if !q.Has("vfs") {
		return false, spelunker_uri, nil
	}

	vfs_url := q.Get("vfs")

	vfs := sqlite3vfshttp.HttpVFS{
		URL: vfs_url,
		RoundTripper: &roundTripper{
			referer:   q.Get("vfs-referer"),
			userAgent: q.Get("vfs-user-agent"),
		},
	}

	err = sqlite3vfs.RegisterVFS("httpvfs", &vfs)

	if err != nil {
		return false, spelunker_uri, fmt.Errorf("Failed to register VFS", "error", err)
	}

	dsn := "spelunker.db?vfs=httpvfs&mode=ro"
	q.Set("dsn", dsn)
	q.Del("vfs")

	u.RawQuery = q.Encode()

	return true, u.String(), nil
}
