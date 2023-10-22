package ordersRepositories

import (
	"encoding/json"
	"fmt"

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
		"o"."id",
		"o"."user_id",
		"o"."transfer_slip",
		(
			SELECT 
				array_to_json(array_agg("pt"))
			FROM (
				SELECT
					"spo"."id",
					"spo"."qty",
					"spo"."product"
				FROM "products_orders" "spo"
				WHERE "spo"."order_id" = "o"."id"

			) AS "pt"
		) AS "products",
		"o"."address",
		"o"."contact",
		"o"."status",
		(
			SELECT 
				SUM(COALESCE(("po"."product"->>'price')::FLOAT*("po"."qty")::FLOAT,0))
			FROM "products_orders" "po"
			WHERE "po"."order_id" = "o"."id"
		) AS "total_paid",
		"o"."created_at",
		"o"."updated_at"
	FROM "orders" "o"
	WHERE "o"."id" = $1

	) AS "t";`

	orderData := &orders.Order{
		TransferSlip: &orders.TransferSlip{},
		Product:      make([]*orders.ProductsOrder, 0),
	}

	rawData := make([]byte, 0)

	if err := r.db.Get(&rawData, query, orderId); err != nil {
		return nil, fmt.Errorf("get order failed:%v", err)
	}

	if err := json.Unmarshal(rawData, &orderData); err != nil {
		return nil, fmt.Errorf("unmarshal order failed:%v", err)
	}

	return orderData, nil
}