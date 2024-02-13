# go-whosonfirst-spelunker-sql

Go package implementing the `whosonfirst/go-whosonfirst-spelunker.Spelunker` interface for use with `database/sql` backed databases.

## Documentation

Documentation is incompete at this time. For starters consult the (also incomplete) documentation in the [whosonfirst/go-whosonfirst-spelunker](https://github.com/whosonfirst/go-whosonfirst-spelunker) package.

## Indexing

For example:

```
$> cd /usr/local/whosonfirst/go-whosonfirst-sqlite-features-index 
$> ./bin/wof-sqlite-index-features-mattn \
	-timings \
	-database-uri mattn:///usr/local/data/ca-search.db \
	-spatial-tables \
	-ancestors \
	-search \
	-concordances \	
	-index-alt-files \
	/usr/local/data/whosonfirst-data-admin-ca
```

## Tools

### server

For example:

```
$> go run -mod readonly -tags "icu json1 fts5" cmd/httpd/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri sql://sqlite3?dsn=file:/usr/local/data/ca-search.db
2024/02/13 08:46:41 INFO Listening for requests address=http://localhost:8080
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-spelunker
* https://github.com/whosonfirst/go-whosonfirst-spelunker-httpd