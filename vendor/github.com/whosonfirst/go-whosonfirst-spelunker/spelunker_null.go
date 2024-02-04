package spelunker

import (
	"context"

	"github.com/aaronland/go-pagination"
)

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

func (s *NullSpelunker) GetDescendants(ctx context.Context, id int64, pg_opts pagination.Options) ([][]byte, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}
