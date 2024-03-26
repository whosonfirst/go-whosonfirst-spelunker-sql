package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/api"
)

func geoJSONHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	opts := &api.GeoJSONHandlerOptions{
		Spelunker: sp,
	}

	return api.GeoJSONHandler(opts)
}

func geoJSONLDHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	opts := &api.GeoJSONLDHandlerOptions{
		Spelunker: sp,
	}

	return api.GeoJSONLDHandler(opts)
}

func sprHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	opts := &api.SPRHandlerOptions{
		Spelunker: sp,
	}

	return api.SPRHandler(opts)
}

func selectHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	// Make this a config/flag
	select_pattern := `properties(?:.[a-zA-Z0-9-_]+){1,}`

	pat, err := regexp.Compile(select_pattern)

	if err != nil {
		slog.Error("Failed to compile select pattern", "pattern", select_pattern, "error", err)
		return nil, fmt.Errorf("Failed to compile select pattern (%s), %w", select_pattern, err)
	}

	opts := &api.SelectHandlerOptions{
		Pattern:   pat,
		Spelunker: sp,
	}

	return api.SelectHandler(opts)
}

func navPlaceHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	opts := &api.NavPlaceHandlerOptions{
		Spelunker:   sp,
		MaxFeatures: 10,
	}

	return api.NavPlaceHandler(opts)
}

func svgHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	sz := api.DefaultSVGSizes()

	opts := &api.SVGHandlerOptions{
		Spelunker: sp,
		Sizes:     sz,
	}

	return api.SVGHandler(opts)
}

func descendantsFacetedHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	opts := &api.DescendantsFacetedHandlerOptions{
		Spelunker: sp,
		// Authenticator: authenticator,
	}

	return api.DescendantsFacetedHandler(opts)
}

func placetypeFacetedHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	opts := &api.PlacetypeFacetedHandlerOptions{
		Spelunker: sp,
		// Authenticator: authenticator,
	}

	return api.PlacetypeFacetedHandler(opts)
}

func recentFacetedHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	opts := &api.RecentFacetedHandlerOptions{
		Spelunker: sp,
		// Authenticator: authenticator,
	}

	return api.RecentFacetedHandler(opts)
}

func hasConcordanceFacetedHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	opts := &api.HasConcordanceFacetedHandlerOptions{
		Spelunker: sp,
		// Authenticator: authenticator,
	}

	return api.HasConcordanceFacetedHandler(opts)
}

func searchFacetedHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	opts := &api.SearchFacetedHandlerOptions{
		Spelunker: sp,
		// Authenticator: authenticator,
	}

	return api.SearchFacetedHandler(opts)
}

func nullIslandFacetedHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	opts := &api.NullIslandFacetedHandlerOptions{
		Spelunker: sp,
		// Authenticator: authenticator,
	}

	return api.NullIslandFacetedHandler(opts)
}
