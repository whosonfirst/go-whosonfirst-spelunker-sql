package spelunker

import (
	"context"
	"time"

	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

// NullSpelunker implements the [Spelunker] interface but returns an `ErrNotImplemented` error for every method.
type NullSpelunker struct {
	Spelunker
}

func init() {
	ctx := context.Background()
	RegisterSpelunker(ctx, "null", NewNullSpelunker)
}

func NewNullSpelunker(ctx context.Context, uri string) (Spelunker, error) {

	s := &NullSpelunker{}

	return s, nil
}

func (s *NullSpelunker) GetById(ctx context.Context, id int64) ([]byte, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) GetDescendants(ctx context.Context, g_opts pagination.Options, id int64, filters ...Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

func (s *NullSpelunker) CountDescendants(ctx context.Context, id int64) (int64, error) {
	return 0, ErrNotImplemented
}

func (s *NullSpelunker) Search(ctx context.Context, pg_opts pagination.Options, q *SearchOptions) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

func (s *NullSpelunker) GetRecent(context.Context, pagination.Options, time.Duration, ...Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}
