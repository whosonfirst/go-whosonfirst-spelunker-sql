CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

GOTAGS=icu json1 fts5

cli:
	go build -mod $(GOMOD) -tags="$(GOTAGS)" -ldflags="$(LDFLAGS)" -o bin/wof-spelunker cmd/wof-spelunker/main.go
	go build -mod $(GOMOD) -tags="$(GOTAGS)" -ldflags="$(LDFLAGS)" -o bin/wof-spelunker-httpd cmd/wof-spelunker-httpd/main.go

SPELUNKER_DATABASE=/usr/local/data/ca-search.db
SPELUNKER_URI=sql://sqlite3?dsn=file:$(SPELUNKER_DATABASE)

server:
	go run -mod $(GOMOD) -tags "$(GOTAGS)" cmd/wof-spelunker-httpd/main.go \
		-server-uri http://localhost:8080 \
		-protomaps-api-key '$(APIKEY)' \
		-spelunker-uri $(SPELUNKER_URI)
