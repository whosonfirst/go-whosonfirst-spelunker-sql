package sql

import (
	"context"
	"fmt"
	"strings"
	"log/slog"
	
	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
	"github.com/whosonfirst/go-whosonfirst-sqlite-spr"
)

func (s *SQLSpelunker) GetAlternatePlacetypes(ctx context.Context) (*spelunker.Faceting, error) {

	facet_counts := make([]*spelunker.FacetCount, 0)

	// TBD alt files...
	q := fmt.Sprintf(`SELECT JSON_EXTRACT(geojson.body, '$.properties.wof:placetype_alt') AS placetype_alt, COUNT(id) AS count FROM %s WHERE placetype_alt != "" AND is_alt=0 GROUP BY placetype_alt ORDER BY count DESC`, tables.GEOJSON_TABLE_NAME)

	slog.Info(q)
	
	rows, err := s.db.QueryContext(ctx, q)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute query, %w", err)
	}

	for rows.Next() {

		var pt string
		var count int64

		err := rows.Scan(&pt, &count)

		if err != nil {
			return nil, fmt.Errorf("Failed to scan row, %w", err)
		}

		f := &spelunker.FacetCount{
			Key:   pt,
			Count: count,
		}

		facet_counts = append(facet_counts, f)
	}

	err = rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close results rows, %w", err)
	}

	f := spelunker.NewFacet("placetype_alt")

	faceting := &spelunker.Faceting{
		Facet:   f,
		Results: facet_counts,
	}

	return faceting, nil
}

func (s *SQLSpelunker) HasAlternatePlacetype(ctx context.Context, pg_opts pagination.Options, pt string, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	where, args, err := s.hasAlternatePlacetypeQueryWhere(pt, filters)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to derive placetype query, %w", err)
	}

	str_where := strings.Join(where, " AND ")
	return s.querySPR(ctx, pg_opts, str_where, args...)
}

func (s *SQLSpelunker) HasAlternatePlacetypeFaceted(ctx context.Context, pt string, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q_where, q_args, err := s.hasAlternatePlacetypeQueryWhere(pt, filters)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive query where statement, %w", err)
	}

	results := make([]*spelunker.Faceting, len(facets))

	// START OF do this in go routines

	for idx, f := range facets {

		q := s.hasPlacetypeQueryFacetStatement(ctx, f, q_where)
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

func (s *SQLSpelunker) hasAlternatePlacetypeQueryWhere(pt string, filters []spelunker.Filter) ([]string, []interface{}, error) {

	where := []string{
		"JSON_EXTRACT(geojson.body, '$.properties.wof:placetype_alt') = ?",
	}

	args := []interface{}{
		pt,
	}

	where, args, err := s.assignFilters(where, args, filters)

	if err != nil {
		return nil, nil, err
	}

	return where, args, nil
}

func (s *SQLSpelunker) hasAlternatePlacetypeQueryFacetStatement(ctx context.Context, facet *spelunker.Facet, where []string) string {

	facet_label := s.facetLabel(facet)

	cols := []string{
		fmt.Sprintf("%s.%s AS %s", tables.GEOJSON_TABLE_NAME, facet_label, facet),
		fmt.Sprintf("COUNT(%s.id) AS count", tables.GEOJSON_TABLE_NAME),
	}

	q := s.hasAlternatePlacetypeQueryStatement(ctx, cols, where)
	return fmt.Sprintf("%s GROUP BY %s.%s ORDER BY count DESC", q, tables.GEOJSON_TABLE_NAME, facet_label)
}

func (s *SQLSpelunker) hasAlternatePlacetypeQueryStatement(ctx context.Context, cols []string, where []string) string {

	str_cols := strings.Join(cols, ",")
	str_where := strings.Join(where, " AND ")

	return fmt.Sprintf("SELECT %s FROM %s WHERE %s", str_cols, tables.GEOJSON_TABLE_NAME, str_where)

}

func (s *SQLSpelunker) queryGeoJSON(ctx context.Context, pg_opts pagination.Options, where string, args ...interface{}) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	if pg_opts != nil {
		limit, offset := s.deriveLimitOffset(pg_opts)
		where = fmt.Sprintf("%s LIMIT %d OFFSET %d", where, limit, offset)
	}

	pg_ch := make(chan pagination.Results)
	results_ch := make(chan wof_spr.StandardPlacesResults)

	done_ch := make(chan bool)
	err_ch := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {

		defer func() {
			done_ch <- true
		}()

		count_q := fmt.Sprintf("SELECT %s.id AS id FROM %s WHERE %s", tables.GEOJSON_TABLE_NAME, tables.GEOJSON_TABLE_NAME, where)
		count, err := s.queryCount(ctx, "id", count_q, args...)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive query count, %w", err)
			return
		}

		var pg_results pagination.Results
		var pg_err error

		if pg_opts != nil {
			pg_results, pg_err = countable.NewResultsFromCountWithOptions(pg_opts, count)
		} else {
			pg_results, pg_err = countable.NewResultsFromCount(count)
		}

		if pg_err != nil {
			err_ch <- fmt.Errorf("Failed to derive pagination results, %w", pg_err)
			return
		}

		pg_ch <- pg_results
	}()

	go func() {

		defer func() {
			done_ch <- true
		}()

		results_q := fmt.Sprintf("SELECT body FROM ? WHERE ?", tables.GEOJSON_TABLE_NAME, where)
		rows, err := s.db.QueryContext(ctx, results_q, args...)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to query where '%s', %w", results_q, err)
			return
		}

		results := make([]wof_spr.StandardPlacesResult, 0)

		for rows.Next() {

			select {
			case <-ctx.Done():
				break
			default:
				// pass
			}

			var body string
			err := rows.Scan(&body)

			if err != nil {
				err_ch <- fmt.Errorf("Failed to read body from row, %w", err)
				return
			}

			spr_row, err := wof_spr.WhosOnFirstSPR([]byte(body))

			if err != nil {
				err_ch <- fmt.Errorf("Failed to create SPR from row, %w", err)
			}

			results = append(results, spr_row)
		}

		err = rows.Close()

		if err != nil {
			err_ch <- fmt.Errorf("Failed to close results rows for descendants, %w", err)
			return
		}

		spr_results := &spr.SQLiteResults{
			Places: results,
		}

		results_ch <- spr_results
	}()

	var pg_results pagination.Results
	var spr_results wof_spr.StandardPlacesResults

	remaining := 2

	for remaining > 0 {
		select {
		case <-done_ch:
			remaining -= 1
		case r := <-pg_ch:
			pg_results = r
		case r := <-results_ch:
			spr_results = r
		case err := <-err_ch:
			return nil, nil, err
		}
	}

	return spr_results, pg_results, nil
}
