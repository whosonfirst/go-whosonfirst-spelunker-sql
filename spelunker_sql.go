package sql

import (
	"context"
	db_sql "database/sql"
	"fmt"
	_ "log/slog"
	"net/url"

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

func (s *SQLSpelunker) GetDescendants(ctx context.Context, id int64) ([][]byte, error) {
	return nil, spelunker.ErrNotImplemented
}
