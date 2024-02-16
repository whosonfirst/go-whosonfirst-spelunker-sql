package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/www"
)

func staticHandlerFunc(ctx context.Context) (http.Handler, error) {

	http_fs := http.FS(run_options.StaticAssets)
	fs_handler := http.FileServer(http_fs)

	return http.StripPrefix(run_options.URIs.Static, fs_handler), nil
}

func indexHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.TemplateHandlerOptions{
		Authenticator: authenticator,
		Templates:     html_templates,
		TemplateName:  "index",
		PageTitle:     "",
		URIs:          uris_table,
	}

	return www.TemplateHandler(opts)
}

func aboutHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.TemplateHandlerOptions{
		Authenticator: authenticator,
		Templates:     html_templates,
		TemplateName:  "about",
		PageTitle:     "About",
		URIs:          uris_table,
	}

	return www.TemplateHandler(opts)
}

func descendantsHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.DescendantsHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.DescendantsHandler(opts)
}

func recentHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.RecentHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.RecentHandler(opts)
}

func idHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.IdHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.IdHandler(opts)
}

func searchHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.SearchHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.SearchHandler(opts)
}
