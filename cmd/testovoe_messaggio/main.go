package main

import (
	"context"
	"os"

	"github.com/ALLGaLL115/testovoe/internal/config"
	"github.com/ALLGaLL115/testovoe/internal/lib/logger/sl"
	"github.com/ALLGaLL115/testovoe/internal/logger"
	"github.com/ALLGaLL115/testovoe/internal/storage/postgresql"
	// "github.com/ALLGaLL115/testovoe-messaggio/internal/config"
	// "github.com/ALLGaLL115/testovoe-messaggio/internal/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)
	log.Debug("debug messages are available")
	log.Info("info messages are available")
	log.Warn("warn messages are available")
	log.Error("error messages are available")

	dbPool, err := postgresql.NewConection(context.TODO(), log, cfg.Database)
	if err != nil {
		log.Error("failed connect to database", sl.Err(err))
		os.Exit(1)
	}

}
