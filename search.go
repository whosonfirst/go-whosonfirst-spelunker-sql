package sql

import (
	"context"
	"strings"

	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
)

func (s *SQLSpelunker) Search(ctx context.Context, pg_opts pagination.Options, search_opts *spelunker.SearchOptions, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	where := []string{
		"names_all MATCH ?",
	}

	str_where := strings.Join(where, " AND ")
	return s.querySearch(ctx, pg_opts, str_where, search_opts.Query)
}

func (s *SQLSpelunker) SearchFaceted(ctx context.Context, search_opts *spelunker.SearchOptions, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	// TO DO
	return nil, spelunker.ErrNotImplemented
}
