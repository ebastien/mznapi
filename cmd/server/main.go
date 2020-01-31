package main

import (
	"github.com/ebastien/mznapi/api"
	log "github.com/sirupsen/logrus"
)

func main() {
	addr := "localhost:8080"
	log.WithField("addr", addr).Info("Starting server")
	srv := api.NewServer(addr, 3)
	srv.Serve()
}
