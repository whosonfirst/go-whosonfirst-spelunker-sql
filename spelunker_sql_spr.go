package sql

import (
	"context"
	"fmt"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
)

func (s *SQLSpelunker) facetSPR(ctx context.Context, facet *spelunker.Facet, where string, args ...interface{}) ([]*spelunker.FacetCount, error) {

	q := fmt.Sprintf("SELECT %s, COUNT(id) AS count FROM %s WHERE %s GROUP BY %s ORDER BY count DESC", facet, tables.SPR_TABLE_NAME, where, facet)
	rows, err := s.db.QueryContext(ctx, q, args...)

	if err != nil {
		return nil, fmt.Errorf("Failed to query facets, %w", err)
	}

	counts := make([]*spelunker.FacetCount, 0)

	for rows.Next() {

		var facet string
		var count int64

		err := rows.Scan(&facet, &count)

		if err != nil {
			return nil, fmt.Errorf("Failed to scan ID, %w", err)
		}

		f := &spelunker.FacetCount{
			Key:   facet,
			Count: count,
		}

		counts = append(counts, f)
	}

	err = rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close results rows for descendants, %w", err)
	}

	return counts, nil
}
