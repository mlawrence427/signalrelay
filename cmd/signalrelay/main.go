package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mlawrence427/signalrelay/internal/server"
	"github.com/mlawrence427/signalrelay/internal/store"
)

func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	stateStore, closeStore, err := openStore()
	if err != nil {
		log.Fatal(err)
	}
	defer closeStore()

	srv := server.New(stateStore)

	log.Printf("signalrelay listening on %s", addr)
	if err := http.ListenAndServe(addr, srv.Routes()); err != nil {
		log.Fatal(err)
	}
}

func openStore() (server.Store, func() error, error) {
	switch strings.ToLower(os.Getenv("SIGNALRELAY_STORE")) {
	case "", "memory":
		return store.NewMemory(), func() error { return nil }, nil
	case "sqlite":
		path := os.Getenv("SIGNALRELAY_DB_PATH")
		if path == "" {
			path = "signalrelay.db"
		}

		sqliteStore, err := store.NewSQLite(path)
		if err != nil {
			return nil, nil, err
		}

		return sqliteStore, sqliteStore.Close, nil
	default:
		return nil, nil, fmt.Errorf("unknown SIGNALRELAY_STORE %q (expected memory or sqlite)", os.Getenv("SIGNALRELAY_STORE"))
	}
}
