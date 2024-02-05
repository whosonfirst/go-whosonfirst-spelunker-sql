CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

SPELUNKER_URI=sql://sqlite3?dsn=file:/usr/local/data/ca-search.db

server:
	go run -mod $(GOMOD) -tags "icu json1 fts5" cmd/httpd/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri $(SPELUNKER_URI)

		# -spelunker-uri $(SPELUNKER_URI) 
server-vfs:
	go run -mod $(GOMOD) -tags "icu json1 fts5" cmd/httpd-vfs/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri 'sql://sqlite3?vfs=http://localhost:8082/test.db'

get_id:
	go run -mod $(GOMOD) cmd/spelunker/main.go \
		-spelunker-uri $(SPELUNKER_URI) \
		-command id \
		-id $(ID)

get_descendants:
	go run -mod $(GOMOD) cmd/spelunker/main.go \
		-spelunker-uri $(SPELUNKER_URI) \
		-command descendants \
		-id $(ID)

search:
	go run -mod $(GOMOD) -tags "icu json1 fts5" cmd/spelunker/main.go \
		-spelunker-uri $(SPELUNKER_URI) \
		-command search \
		-query $(NAME)
