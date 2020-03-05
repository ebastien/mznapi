.DEFAULT_GOAL: build

.PHONY: build
build: build/server

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm build/*

.PHONE: run
run:
	go run ./cmd/server/main.go

build/%: cmd/%/main.go
	go build -o $@ $< 
