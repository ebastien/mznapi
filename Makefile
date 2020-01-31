.PHONY: build
build: build/server

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm build/*

build/%: cmd/%/main.go
	go build -o $@ $< 
