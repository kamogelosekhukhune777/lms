package orderdb

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kamogelosekhukhune777/lms/business/domain/orderbus"
	"github.com/kamogelosekhukhune777/lms/business/sdk/sqldb"
	"github.com/kamogelosekhukhune777/lms/foundation/logger"
)

// Store manages the set of APIs for product database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (orderbus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

func (s *Store) Create(ctx context.Context, ord orderbus.Order) error {
	const q = ``

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBOrder(ord)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) QueryByID(ctx context.Context, orderID uuid.UUID) (orderbus.Order, error) {
	data := struct {
		ID string `db:"order_id"`
	}{
		ID: orderID.String(),
	}

	const q = ``

	var dbOrd order
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbOrd); err != nil {
		return orderbus.Order{}, fmt.Errorf("db: %w", err)
	}

	return toBusOrder(dbOrd)
}

func (s *Store) Update(ctx context.Context, ord orderbus.Order) error {
	const q = ``

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBOrder(ord)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}
