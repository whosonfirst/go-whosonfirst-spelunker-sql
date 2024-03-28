package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
)

func (s *SQLSpelunker) VisitingNullIsland(ctx context.Context, pg_opts pagination.Options, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	where, args, err := s.visitingNullIslandQueryWhere(filters)

	if err != nil {
		return nil, nil, err
	}

	str_where := strings.Join(where, " AND ")
	return s.querySPR(ctx, pg_opts, str_where, args...)
}

func (s *SQLSpelunker) VisitingNullIslandFaceted(ctx context.Context, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q_where, q_args, err := s.visitingNullIslandQueryWhere(filters)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive query where statement, %w", err)
	}

	results := make([]*spelunker.Faceting, len(facets))

	// START OF do this in go routines

	for idx, f := range facets {

		q := s.visitingNullIslandQueryFacetStatement(ctx, f, q_where)
		// slog.Info("FACET", "q", q, "args", q_args)

		counts, err := s.facetWithQuery(ctx, q, q_args...)

		if err != nil {
			return nil, fmt.Errorf("Failed to facet columns, %w", err)
		}

		fc := &spelunker.Faceting{
			Facet:   f,
			Results: counts,
		}

		results[idx] = fc
	}

	// END OF do this in go routines

	return results, nil

}

func (s *SQLSpelunker) visitingNullIslandQueryWhere(filters []spelunker.Filter) ([]string, []interface{}, error) {

	where := []string{
		"latitude = ?",
		"longitude = ?",
	}

	args := []interface{}{
		0.0,
		0.0,
	}

	for _, f := range filters {

		switch f.Scheme() {
		case spelunker.COUNTRY_FILTER_SCHEME:
			where = append(where, "country = ?")
			args = append(args, f.Value())
		case spelunker.PLACETYPE_FILTER_SCHEME:
			where = append(where, "placetype = ?")
			args = append(args, f.Value())
		case spelunker.IS_CURRENT_FILTER_SCHEME:
			where = append(where, "is_current = ?")
			args = append(args, f.Value())
		default:
			return nil, nil, fmt.Errorf("Invalid or unsupported filter scheme, %s", f.Scheme())
		}
	}

	return where, args, nil
}

func (s *SQLSpelunker) visitingNullIslandQueryFacetStatement(ctx context.Context, facet *spelunker.Facet, where []string) string {

	var facet_label string

	switch facet.Property {
	case "iscurrent":
		facet_label = "is_current"
	default:
		facet_label = facet.Property
	}

	cols := []string{
		fmt.Sprintf("%s.%s AS %s", tables.SPR_TABLE_NAME, facet_label, facet),
		fmt.Sprintf("COUNT(%s.id) AS count", tables.SPR_TABLE_NAME),
	}

	q := s.visitingNullIslandQueryStatement(ctx, cols, where)
	return fmt.Sprintf("%s GROUP BY %s.%s ORDER BY count DESC", q, tables.SPR_TABLE_NAME, facet_label)
}

func (s *SQLSpelunker) visitingNullIslandQueryStatement(ctx context.Context, cols []string, where []string) string {

	str_cols := strings.Join(cols, ",")
	str_where := strings.Join(where, " AND ")

	return fmt.Sprintf("SELECT %s FROM %s WHERE %s", str_cols, tables.SPR_TABLE_NAME, str_where)

}
