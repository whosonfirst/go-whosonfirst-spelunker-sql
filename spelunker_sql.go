package sql

import (
	"context"
	db_sql "database/sql"
	"fmt"
	_ "log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
	"github.com/whosonfirst/go-whosonfirst-sqlite-spr"
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

	q := fmt.Sprintf("SELECT body FROM %s WHERE id = ?", tables.GEOJSON_TABLE_NAME)
	return s.getById(ctx, q, id)
}

func (s *SQLSpelunker) GetAlternateGeometryById(ctx context.Context, id int64, alt_geom *uri.AltGeom) ([]byte, error) {

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

func (s *SQLSpelunker) GetDescendants(ctx context.Context, pg_opts pagination.Options, id int64, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	where := []string{
		"instr(belongsto, ?) > 0",
	}

	args := []interface{}{
		id,
	}

	for _, f := range filters {

		switch f.Scheme() {
		case spelunker.COUNTRY_FILTER_SCHEME:
			where = append(where, "country = ?")
			args = append(args, f.Value())
		case spelunker.PLACETYPE_FILTER_SCHEME:
			where = append(where, "placetype = ?")
			args = append(args, f.Value())
		default:
			return nil, nil, fmt.Errorf("Invalid or unsupported filter scheme, %s", f.Scheme())
		}

	}

	str_where := strings.Join(where, " AND ")

	return s.querySPR(ctx, pg_opts, str_where, args...)
}

func (s *SQLSpelunker) GetDescendantsFaceted(ctx context.Context, id int64, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	where := []string{
		"instr(belongsto, ?) > 0",
	}

	args := []interface{}{
		id,
	}

	for _, f := range filters {

		switch f.Scheme() {
		case spelunker.COUNTRY_FILTER_SCHEME:
			where = append(where, "country = ?")
			args = append(args, f.Value())
		case spelunker.PLACETYPE_FILTER_SCHEME:
			where = append(where, "placetype = ?")
			args = append(args, f.Value())
		default:
			return nil, fmt.Errorf("Invalid or unsupported filter scheme, %s", f.Scheme())
		}

	}

	str_where := strings.Join(where, " AND ")

	// START OF do this in go routines

	f := facets[0]

	counts, err := s.facetSPR(ctx, f, str_where, args...)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive facets for %s, %w", f, err)
	}

	results := []*spelunker.Faceting{
		&spelunker.Faceting{
			Facet:   f,
			Results: counts,
		},
	}

	// END OF do this in go routines

	return results, nil
}

func (s *SQLSpelunker) CountDescendants(ctx context.Context, id int64) (int64, error) {

	var count int64

	q := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE instr(belongsto, ?)", tables.SPR_TABLE_NAME)

	row := s.db.QueryRowContext(ctx, q, id)
	err := row.Scan(&count)

	switch {
	case err == db_sql.ErrNoRows:
		return 0, spelunker.ErrNotFound
	case err != nil:
		return 0, fmt.Errorf("Failed to execute count descendants query for %d, %w", id, err)
	default:
		return count, nil
	}
}

func (s *SQLSpelunker) HasPlacetype(ctx context.Context, pg_opts pagination.Options, pt *placetypes.WOFPlacetype, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	where := []string{
		"placetype = ?",
	}

	args := []interface{}{
		pt.Name,
	}

	for _, f := range filters {

		switch f.Scheme() {
		case spelunker.COUNTRY_FILTER_SCHEME:
			where = append(where, "country = ?")
			args = append(args, f.Value())
		case spelunker.PLACETYPE_FILTER_SCHEME:
			where = append(where, "placetype = ?")
			args = append(args, f.Value())
		default:
			return nil, nil, fmt.Errorf("Invalid or unsupported filter scheme, %s", f.Scheme())
		}

	}

	str_where := strings.Join(where, " AND ")

	return s.querySPR(ctx, pg_opts, str_where, args...)
}

func (s *SQLSpelunker) Search(ctx context.Context, pg_opts pagination.Options, search_opts *spelunker.SearchOptions) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	where := "names_all MATCH ?"
	return s.querySearch(ctx, pg_opts, where, search_opts.Query)
}

func (s *SQLSpelunker) GetRecent(ctx context.Context, pg_opts pagination.Options, d time.Duration, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	now := time.Now()
	then := now.Unix() - int64(d.Seconds())

	where := "lastmodified >= ? ORDER BY lastmodified DESC"
	return s.querySPR(ctx, pg_opts, where, then)
}

func (s *SQLSpelunker) GetPlacetypes(ctx context.Context) (*spelunker.Faceting, error) {

	facet_counts := make([]*spelunker.FacetCount, 0)

	q := fmt.Sprintf("SELECT placetype, COUNT(id) AS count FROM %s GROUP BY placetype ORDER BY count DESC", tables.SPR_TABLE_NAME)

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

	f := spelunker.NewFacet("placetype")

	faceting := &spelunker.Faceting{
		Facet:   f,
		Results: facet_counts,
	}

	return faceting, nil
}

func (s *SQLSpelunker) GetConcordances(ctx context.Context) (*spelunker.Faceting, error) {

	facet_counts := make([]*spelunker.FacetCount, 0)

	q := fmt.Sprintf("SELECT other_source, COUNT(other_id) AS count FROM %s GROUP BY other_source ORDER BY count DESC", tables.CONCORDANCES_TABLE_NAME)

	rows, err := s.db.QueryContext(ctx, q)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute query, %w", err)
	}

	for rows.Next() {

		var source string
		var count int64

		err := rows.Scan(&source, &count)

		if err != nil {
			return nil, fmt.Errorf("Failed to scan row, %w", err)
		}

		nspred := strings.Split(source, ":")
		ns := nspred[0]

		f := &spelunker.FacetCount{
			Key:   ns,
			Count: count,
		}

		facet_counts = append(facet_counts, f)
	}

	err = rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close results rows, %w", err)
	}

	f := spelunker.NewFacet("concordance")

	faceting := &spelunker.Faceting{
		Facet:   f,
		Results: facet_counts,
	}

	return faceting, nil
}

func (s *SQLSpelunker) HasConcordance(ctx context.Context, pg_opts pagination.Options, namespace string, predicate string, value string, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	where := make([]string, 0)
	args := make([]interface{}, 0)

	switch {
	case namespace != "" && predicate != "":
		where = append(where, "other_source = ?")
		args = append(args, fmt.Sprintf("%s:%s", namespace, predicate))
	case namespace != "":
		where = append(where, "other_source LIKE ?")
		args = append(args, namespace+":%")
	case predicate != "":
		where = append(where, "other_source LIKE ?")
		args = append(args, "%:"+predicate)
	default:
		return nil, nil, fmt.Errorf("Missing namespace and predicate")
	}

	if value != "" {
		where = append(where, "other_id = ?")
		args = append(args, value)
	}

	str_where := strings.Join(where, " AND ")

	q := fmt.Sprintf("SELECT id FROM %s WHERE %s", tables.CONCORDANCES_TABLE_NAME, str_where)

	rows, err := s.db.QueryContext(ctx, q, args...)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to execute query, %w", err)
	}

	ids := make([]interface{}, 0)
	qms := make([]string, 0)

	for rows.Next() {

		var str_id int64

		err := rows.Scan(&str_id)

		if err != nil {
			return nil, nil, fmt.Errorf("Failed to scan row, %w", err)
		}

		ids = append(ids, str_id)
		qms = append(qms, "?")
	}

	err = rows.Close()

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to close results rows, %w", err)
	}

	if len(ids) == 0 {

		var pg_results pagination.Results
		var pg_err error

		if pg_opts != nil {
			pg_results, pg_err = countable.NewResultsFromCountWithOptions(pg_opts, 0)
		} else {
			pg_results, pg_err = countable.NewResultsFromCount(0)
		}

		if pg_err != nil {
			return nil, nil, fmt.Errorf("Failed to create pagination results, %w", err)
		}

		results := make([]wof_spr.StandardPlacesResult, 0)

		spr_results := &spr.SQLiteResults{
			Places: results,
		}

		return spr_results, pg_results, nil
	}

	spr_where := fmt.Sprintf("id IN (%s)", strings.Join(qms, ","))
	return s.querySPR(ctx, pg_opts, spr_where, ids...)
}
