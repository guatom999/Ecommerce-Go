package ordersRepositories

import (
	"github.com/guatom999/Ecommerce-Go/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IOrderRepository interface {
	FindOneProduct(orderId string) (*orders.Order, error)
}

type orderRepository struct {
	db *sqlx.DB
}

func OrderRepository(db *sqlx.DB) IOrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindOneProduct(orderId string) (*orders.Order, error) {

	query := `
	SELECT 
		to_jsonb("t")
	FROM (
		SELECT 
		FROM "orders" "o"
		WHERE "o"."id" = $1

	) AS "t";`

	return nil, nil
}
