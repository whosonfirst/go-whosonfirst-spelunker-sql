package sql

import (
	"context"
	"fmt"
	_ "log/slog"
	"math"
	"strings"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
	"github.com/whosonfirst/go-whosonfirst-sqlite-spr"
)

func (s *SQLSpelunker) queryCount(ctx context.Context, col string, q string, args ...interface{}) (int64, error) {

	parts := strings.Split(q, " FROM ")
	parts = strings.Split(parts[1], " LIMIT ")
	parts = strings.Split(parts[0], " ORDER ")

	conditions := parts[0]

	count_query := fmt.Sprintf("SELECT COUNT(%s) FROM %s", col, conditions)

	row := s.db.QueryRowContext(ctx, count_query, args...)

	var count int64
	err := row.Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("Failed to execute count query '%s', %w", count_query, err)
	}

	return count, nil
}

func (s *SQLSpelunker) deriveLimitOffset(pg_opts pagination.Options) (int, int) {

	page_num := countable.PageFromOptions(pg_opts)
	page := int(math.Max(1.0, float64(page_num)))

	per_page := int(math.Max(1.0, float64(pg_opts.PerPage())))
	spill := int(math.Max(1.0, float64(pg_opts.Spill())))

	if spill >= per_page {
		spill = per_page - 1
	}

	offset := 0
	limit := per_page

	offset = (page - 1) * per_page

	return limit, offset
}

func (s *SQLSpelunker) selectSPR(ctx context.Context, where string) string {
	return fmt.Sprintf(`SELECT 
		id, parent_id, name, placetype,
		inception, cessation,
		country, repo,
		latitude, longitude,
		min_latitude, min_longitude,
		max_latitude, max_longitude,
		is_current, is_deprecated, is_ceased,is_superseded, is_superseding,
		supersedes, superseded_by, belongsto,
		is_alt, alt_label,
		lastmodified
	FROM %s WHERE %s`, tables.SPR_TABLE_NAME, where)
}

func (s *SQLSpelunker) querySPR(ctx context.Context, pg_opts pagination.Options, where string, args ...interface{}) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	limit, offset := s.deriveLimitOffset(pg_opts)

	where = fmt.Sprintf("%s LIMIT %d OFFSET %d", where, limit, offset)

	pg_ch := make(chan pagination.Results)
	results_ch := make(chan wof_spr.StandardPlacesResults)
	err_ch := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {

		count_q := fmt.Sprintf("SELECT id FROM %s WHERE %s", tables.SPR_TABLE_NAME, where)

		count, err := s.queryCount(ctx, "id", count_q, args...)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive query count, %w", err)
			return
		}

		pg_results, err := countable.NewResultsFromCountWithOptions(pg_opts, count)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive pagination results, %w", err)
			return
		}

		pg_ch <- pg_results
	}()

	go func() {

		results_q := s.selectSPR(ctx, where)

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

			spr_row, err := spr.RetrieveSPRWithRows(ctx, rows)

			if err != nil {
				err_ch <- fmt.Errorf("Failed to derive SPR from row, %w", err)
				return
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
		case r := <-pg_ch:
			pg_results = r
			remaining -= 1
		case r := <-results_ch:
			spr_results = r
			remaining -= 1
		case err := <-err_ch:
			return nil, nil, err
		}
	}

	return spr_results, pg_results, nil
}
