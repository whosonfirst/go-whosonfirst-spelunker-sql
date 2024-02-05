package sql

import (
	"context"
	db_sql "database/sql"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
)

type SQLSpelunker struct {
	spelunker.Spelunker
	engine string
	db     *db_sql.DB
}

func init() {
	ctx := context.Background()
	spelunker.RegisterSpelunker(ctx, "sql", NewSQLSpelunker)
}

func NewSQLSpelunker(ctx context.Context, uri string) (spelunker.Spelunker, error) {

	slog.Info("NEW", "uri", uri)
	
	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	engine := u.Host

	q := u.Query()

	dsn := q.Get("dsn")

	slog.Info("CONNECT", "engine", engine, "dsn", dsn)
	db, err := db_sql.Open(engine, dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection, %w", err)
	}

	s := &SQLSpelunker{
		engine: engine,
		db:     db,
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

func (s *SQLSpelunker) GetDescendants(ctx context.Context, id int64, pg_opts pagination.Options) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	where := "instr(belongsto, ?) > 0"
	return s.querySPR(ctx, pg_opts, where, id)
}

func (s *SQLSpelunker) Search(ctx context.Context, search_opts *spelunker.SearchOptions, pg_opts pagination.Options) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	where := "names_all MATCH ?"
	return s.querySearch(ctx, pg_opts, where, search_opts.Query)
}
