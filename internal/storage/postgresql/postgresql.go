package postgresql

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ALLGaLL115/testovoe-messaggio/internal/config"
	"github.com/ALLGaLL115/testovoe-messaggio/lib/logger/sl"
	"github.com/ALLGaLL115/testovoe-messaggio/lib/storage/repetable"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(ctx context.Context, log *slog.Logger, cfg config.DataBase) (*pgxpool.Pool, error) {
	const op = "database.postgresql.NewClient"

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	var pool *pgxpool.Pool

	err := repetable.DoWithTries(func() error {
		log.Info("Database connection attemp")
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)

		defer cancel()

		pool, _ = pgxpool.New(ctx, dsn)
		err := pool.Ping(ctx)

		if err != nil {
			log.Error("database connection error")
		}

		return err
	}, cfg.Attempts, cfg.Delay)

	if err != nil {
		log.Error("falied connect to database", sl.OpError(op, err))
		return nil, err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		stopSignal := <-stop
		log.Info("stoppping database connection", slog.String("op", op), slog.String("signal", stopSignal.String()))
		pool.Close()
		log.Info("database was stopped")
	}()

	return pool, nil
}
