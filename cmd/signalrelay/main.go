package main

import (
	"log"
	"net/http"

	"github.com/mlawrence427/signalrelay/internal/config"
	"github.com/mlawrence427/signalrelay/internal/server"
	"github.com/mlawrence427/signalrelay/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	stateStore, closeStore, err := openStore(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer closeStore()

	srv := server.New(stateStore)

	logStartup(cfg)
	if err := http.ListenAndServe(cfg.Addr, srv.Routes()); err != nil {
		log.Fatal(err)
	}
}

func openStore(cfg config.Config) (server.Store, func() error, error) {
	switch cfg.Store {
	case config.StoreMemory:
		return store.NewMemory(), func() error { return nil }, nil
	case config.StoreSQLite:
		sqliteStore, err := store.NewSQLite(cfg.DBPath)
		if err != nil {
			return nil, nil, err
		}

		return sqliteStore, sqliteStore.Close, nil
	default:
		return nil, nil, config.StoreError(cfg.Store)
	}
}

func logStartup(cfg config.Config) {
	if cfg.Store == config.StoreSQLite {
		log.Printf("signalrelay starting addr=%s store=%s db_path=%s", cfg.Addr, cfg.Store, cfg.DBPath)
		return
	}

	log.Printf("signalrelay starting addr=%s store=%s", cfg.Addr, cfg.Store)
}
