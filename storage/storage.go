package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"zero/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	CreateOrder(ctx context.Context, orderdb *model.OrderDb) error
	GetOrderDb(ctx context.Context, orderUid string) (model.OrderDb, error)
	GetOrdersDb(ctx context.Context) ([]model.OrderDb, error)
}

type postgres struct {
	db *pgxpool.Pool
}

func NewPgStorage(ctx context.Context, connString string) (*postgres, error) {
	var pgInstance *postgres
	var pgOnce sync.Once
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, connString)
		if err != nil {
			fmt.Printf("Unable to create connection pool: %v\n", err)
			return
		}

		pgInstance = &postgres{db}
	})

	return pgInstance, nil
}

func (pg *postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *postgres) Close() {
	pg.db.Close()
}

func (pg *postgres) CreateOrder(ctx context.Context, orderdb *model.OrderDb) error {

	query := `INSERT INTO orders (order_uid, data_json) 
	VALUES (@order_uid, @data_json)`
	args := pgx.NamedArgs{
		"order_uid": orderdb.OrderUID,
		"data_json": orderdb.DataJson,
	}
	_, err := pg.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	return nil
}

func (pg *postgres) GetOrderDb(ctx context.Context, orderUid string) (model.OrderDb, error) {
	query := fmt.Sprintf(`SELECT order_uid, data_json FROM orders WHERE order_uid = '%s'`, orderUid)

	rows, err := pg.db.Query(ctx, query)

	var orderDb model.OrderDb

	if err != nil {
		return orderDb, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	orderDb, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[model.OrderDb])

	if errors.Is(err, pgx.ErrNoRows) {
		return orderDb, errors.New("username not found")
	}

	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return orderDb, err
	}

	return orderDb, nil
}

func (pg *postgres) GetOrdersDb(ctx context.Context) ([]model.OrderDb, error) {
	query := `SELECT order_uid, data_json FROM orders`

	rows, err := pg.db.Query(ctx, query)

	var orderDb []model.OrderDb

	if err != nil {
		return orderDb, fmt.Errorf("unable to query: %w", err)
	}
	defer rows.Close()

	orderDb, err = pgx.CollectRows(rows, pgx.RowToStructByName[model.OrderDb])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return orderDb, err
	}

	return orderDb, nil
}
