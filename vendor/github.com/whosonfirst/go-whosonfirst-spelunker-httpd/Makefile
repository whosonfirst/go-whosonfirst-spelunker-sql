CWD=$(shell pwd)

GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

debug:
	go run -mod $(GOMOD) cmd/server/main.go \
		-server-uri http://localhost:8080 \
		-spelunker-uri null://
