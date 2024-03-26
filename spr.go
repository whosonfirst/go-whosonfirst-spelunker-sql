package sql

// Common code for querying the `spr` table.

import (
	"context"
	"fmt"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
)

func (s *SQLSpelunker) facetSPR(ctx context.Context, facet *spelunker.Facet, where string, args ...interface{}) ([]*spelunker.FacetCount, error) {

	q := fmt.Sprintf("SELECT %s, COUNT(id) AS count FROM %s WHERE %s GROUP BY %s ORDER BY count DESC", facet, tables.SPR_TABLE_NAME, where, facet)

	return s.facetWithQuery(ctx, q, args...)
}
