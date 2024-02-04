package sql

import (
	"context"
	db_sql "database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
)

type SQLSpelunker struct {
	spelunker.Spelunker
	db *db_sql.DB
}

func init() {
	ctx := context.Background()
	spelunker.RegisterSpelunker(ctx, "sql", NewSQLSpelunker)
}

func NewSQLSpelunker(ctx context.Context, uri string) (spelunker.Spelunker, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	engine := u.Host

	q := u.Query()

	dsn := q.Get("dsn")

	db, err := db_sql.Open(engine, dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection, %w", err)
	}

	s := &SQLSpelunker{
		db: db,
	}

	return s, nil
}

func (s *SQLSpelunker) GetById(ctx context.Context, id int64) ([]byte, error) {

	var body []byte

	q := fmt.Sprintf("SELECT body FROM %s WHERE id = ?", tables.GEOJSON_TABLE_NAME)

	rsp := s.db.QueryRowContext(ctx, q, id)

	err := rsp.Scan(&body)

	switch {
	case err == db_sql.ErrNoRows:
		return nil, spelunker.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("Failed to execute get by id query for %d, %w", id, err)
	default:
		return body, nil
	}
}

func (s *SQLSpelunker) GetDescendants(ctx context.Context, id int64, pg_opts pagination.Options) ([][]byte, pagination.Results, error) {

	limit := pg_opts.PerPage
	offset := countable.PageFromOptions(pg_opts)

	// This is probably SQLite specific...
	q := fmt.Sprintf("SELECT id FROM %s WHERE instr(belongsto, ?) > 0 LIMIT %d, OFFSET %d", tables.SPR_TABLE_NAME, limit, offset)

	slog.Info("GetDescendants", "query", q)

	count, err := s.queryCount(ctx, "id", q, id)

	if err != nil {
		return nil, nil, err
	}

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to derive pagination results, %w", err)
	}

	pg_results, err := countable.NewResultsFromCountWithOptions(pg_opts, count)

	q = fmt.Sprintf("SELECT g.body FROM %s g, %s s WHERE g.id=s.id AND instr(s.belongsto, ?) > 0 LIMIT %d OFFSET %d", tables.GEOJSON_TABLE_NAME, tables.SPR_TABLE_NAME, limit, offset)

	slog.Info("GetDescendants", "query", q)

	rows, err := s.db.QueryContext(ctx, q, id)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to query descendants, %w", err)
	}

	results := make([][]byte, 0)

	for rows.Next() {

		var body []byte

		err := rows.Scan(&body)

		if err != nil {
			return nil, nil, fmt.Errorf("Failed to scan descendants row, %w", err)
		}

		results = append(results, body)
	}

	err = rows.Close()

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to close results rows for descendants, %w", err)
	}

	return results, pg_results, nil
}

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
