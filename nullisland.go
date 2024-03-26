package sql

import (
	"context"

	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	_ "github.com/whosonfirst/go-whosonfirst-sql/tables"
)

func (s *SQLSpelunker) VisitingNullIsland(ctx context.Context, pg_opts pagination.Options, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	// TO DO
	return nil, nil, spelunker.ErrNotImplemented
}

func (s *SQLSpelunker) VisitingNullIslandFaceted(ctx context.Context, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	// TO DO
	return nil, spelunker.ErrNotImplemented
}
