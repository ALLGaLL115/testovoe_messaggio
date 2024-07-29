package postgresql

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ALLGaLL115/testovoe-messaggio/internal/config"

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

}
