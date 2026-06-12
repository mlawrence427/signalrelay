package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mlawrence427/signalrelay/internal/server"
	"github.com/mlawrence427/signalrelay/internal/store"
)

func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := server.New(store.NewMemory())

	log.Printf("signalrelay listening on %s", addr)
	if err := http.ListenAndServe(addr, srv.Routes()); err != nil {
		log.Fatal(err)
	}
}
