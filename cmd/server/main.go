package main

import (
	"flag"

	"github.com/ebastien/mznapi/api"
	"github.com/ebastien/mznapi/store"
	log "github.com/sirupsen/logrus"
)

func main() {

	addr := flag.String("addr", "localhost:8080", "address to listen to")
	debug := flag.Bool("debug", false, "enable debug logs")

	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.WithField("addr", *addr).Info("Starting server")
	srv := api.NewServer(*addr, 3, store.NewMemoryStore())
	srv.Serve()
}
