package sql

import (
	"context"
	db_sql "database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker/document"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
	"github.com/whosonfirst/go-whosonfirst-uri"
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

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	engine := u.Host

	q := u.Query()

	dsn := q.Get("dsn")

	if dsn == "" {
		return nil, fmt.Errorf("Missing ?dsn= parameter")
	}

	slog.Info("DSN", "dsn", dsn)

	db, err := db_sql.Open(engine, dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection, %w", err)
	}

	db.SetMaxOpenConns(1)

	s := &SQLSpelunker{
		engine: engine,
		db:     db,
	}

	return s, nil
}

func (s *SQLSpelunker) GetRecordForId(ctx context.Context, id int64) ([]byte, error) {

	// TBD - replace this with a dedicated "spelunker" table
	// https://github.com/whosonfirst/go-whosonfirst-sql/blob/spelunker/tables/spelunker.sqlite.schema

	q := fmt.Sprintf("SELECT body FROM %s WHERE id = ?", tables.GEOJSON_TABLE_NAME)
	body, err := s.getById(ctx, q, id)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve record, %w", err)
	}

	return document.PrepareSpelunkerV2Document(ctx, body)
}

func (s *SQLSpelunker) GetFeatureForId(ctx context.Context, id int64, uri_args *uri.URIArgs) ([]byte, error) {

	alt_geom := uri_args.AltGeom
	alt_label, err := alt_geom.String()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive label from alt geom, %w", err)
	}

	q := fmt.Sprintf("SELECT body FROM %s WHERE id = ? AND alt_label = ?", tables.GEOJSON_TABLE_NAME)
	return s.getById(ctx, q, id, alt_label)
}

func (s *SQLSpelunker) getById(ctx context.Context, q string, args ...interface{}) ([]byte, error) {

	var body []byte

	rsp := s.db.QueryRowContext(ctx, q, args...)

	err := rsp.Scan(&body)

	switch {
	case err == db_sql.ErrNoRows:
		return nil, spelunker.ErrNotFound
	case err != nil:
		return nil, fmt.Errorf("Failed to execute get by id query, %w", err)
	default:
		return body, nil
	}
}

func (s *SQLSpelunker) Search(ctx context.Context, pg_opts pagination.Options, search_opts *spelunker.SearchOptions, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	where := []string{
		"names_all MATCH ?",
	}

	str_where := strings.Join(where, " AND ")
	return s.querySearch(ctx, pg_opts, str_where, search_opts.Query)
}

func (s *SQLSpelunker) GetRecent(ctx context.Context, pg_opts pagination.Options, d time.Duration, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	now := time.Now()
	then := now.Unix() - int64(d.Seconds())

	where := []string{
		"lastmodified >= ? ORDER BY lastmodified DESC",
	}

	str_where := strings.Join(where, " AND ")
	return s.querySPR(ctx, pg_opts, str_where, then)
}
