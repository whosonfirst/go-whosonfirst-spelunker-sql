package sql

import (
	"context"
	db_sql "database/sql"
	"fmt"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker/document"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

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
