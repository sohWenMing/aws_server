package dbpostgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStringer interface {
	GetDBString() string
}

func ConnectToDB(s DBStringer, ctx context.Context) (pool *pgxpool.Pool, err error) {
	dbPool, err := pgxpool.New(ctx, s.GetDBString())
	if err != nil {
		return nil, err
	}
	return dbPool, nil
}
