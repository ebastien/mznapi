package main

import (
	"github.com/ebastien/mznapi/api"
)

func main() {
	srv := api.NewServer("localhost:8080", 3)
	srv.Serve()
}
