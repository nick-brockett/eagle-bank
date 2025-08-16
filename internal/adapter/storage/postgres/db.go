package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type DBContext struct {
	DB *sqlx.DB
}

// NewDBContext returns a DBContext
func NewDBContext(ctx context.Context, DBConfig Config) (*DBContext, error) {
	db, err := OpenDB(ctx, DBConfig)
	if err != nil {
		return nil, err
	}
	return &DBContext{
		DB: db,
	}, nil
}

// OpenDB returns a PostgresSQL sqlx.DB.
func OpenDB(_ context.Context, dbCfg Config) (*sqlx.DB, error) {

	dataSourceName, err := dbCfg.PostgresConnString()
	if err != nil {
		return nil, err
	}
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to postgres")
	}

	err = db.Ping()
	
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping postgres")
	} else {
		fmt.Printf("Successfully pinged postgres on host %s \n", dbCfg.Host)
	}
	return db, nil
}

func (ctx *DBContext) Close() error {
	err := ctx.DB.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close DB")
	}
	return nil
}
