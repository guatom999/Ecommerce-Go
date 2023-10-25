package ordersPatterns

import (
	"context"
	"fmt"
	"time"

	"github.com/guatom999/Ecommerce-Go/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IInsertOrderBuilder interface {
	initTransaction() error
	insertOrder() error
	insertProductsOrders() error
	commit() error
	getOrderId() string
}

type insertOrderBuilder struct {
	db  *sqlx.DB
	req *orders.Order
	tx  *sqlx.Tx
}

func InsertOrderBuilder(db *sqlx.DB, req *orders.Order) IInsertOrderBuilder {
	return &insertOrderBuilder{
		db:  db,
		req: req,
	}
}

func (b *insertOrderBuilder) initTransaction() error {

	tx, err := b.db.BeginTxx(context.Background(), nil)

	if err != nil {
		return err
	}

	b.tx = tx

	return nil
}

func (b *insertOrderBuilder) insertOrder() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	query := `
	INSERT INTO "orders" ( 
		"user_id",
		"contact",
		"address",
		"tranfer_slip",
		"status"
		)
	VALUES
		( $1 ,$2 ,$3 ,$4 ,$5 ) RETURNING "id";`

	if err := b.tx.QueryRowxContext(
		ctx, query,
		b.req.UserId,
		b.req.Contact,
		b.req.Address,
		b.req.TransferSlip,
		b.req.Status,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert order failed:%v", err)
	}

	return nil
}
func (b *insertOrderBuilder) insertProductsOrders() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	query := `
	INSERT INTO "products_orders" (
		"order_id",
		"qty",
		"product"
	) VALUES`

	values := make([]any, 0)
	lastIndex := 0

	for i := range b.req.Product {
		values = append(
			values,
			b.req.Id,
			b.req.Product[i].Quantity,
			b.req.Product[i],
		)

		if len(b.req.Product)-1 != i {
			query += fmt.Sprintf(`
			( $%d , $%d , $%d ),`, lastIndex+1, lastIndex+2, lastIndex+3)
		} else {
			query += fmt.Sprintf(`
			( $%d , $%d , $%d );`, lastIndex+1, lastIndex+2, lastIndex+3)
		}

		lastIndex += 3
	}

	if _, err := b.tx.ExecContext(ctx, query, values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert product_orders failed:%v", err)
	}

	return nil
}

func (b *insertOrderBuilder) commit() error {

	if err := b.tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (b *insertOrderBuilder) getOrderId() string {

	return b.req.Id
}

type IInsertOrderEngineer interface {
	InsertOrder() (string, error)
}

type insertOrderEngineer struct {
	builder IInsertOrderBuilder
}

func InsertOrderEngineer(builder IInsertOrderBuilder) *insertOrderEngineer {
	return &insertOrderEngineer{builder: builder}
}

func (en *insertOrderEngineer) Chon() string {
	return "chon"
}

func (en *insertOrderEngineer) InsertOrder() (string, error) {

	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}
	if err := en.builder.insertOrder(); err != nil {
		return "", err
	}
	if err := en.builder.insertProductsOrders(); err != nil {
		return "", err
	}
	if err := en.builder.commit(); err != nil {
		return "", err
	}

	return en.builder.getOrderId(), nil
}
