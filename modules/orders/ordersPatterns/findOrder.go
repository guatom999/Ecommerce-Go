package ordersPatterns

import (
	"fmt"
	"strings"

	"github.com/guatom999/Ecommerce-Go/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IFindOrderBuilder interface {
	initQuery()
	initCountQuery()
	buildWhereSearch()
	buildWhereStatus()
	buildWhereDate()
	builSort()
	buildPaginate()
	FinalQuery()
	closeQuery()
	getQuery() string
	setQuery(query string)
	getValues() []any
	setValues(data []any)
	setLastIndex(n int)
	getDb() *sqlx.DB
	reset()
}

type findOrderBuilder struct {
	db        *sqlx.DB
	req       *orders.OrderFilter
	query     string
	values    []any
	lastIndex int
}

func FindOrderBuilder(db *sqlx.DB, req *orders.OrderFilter) IFindOrderBuilder {
	return &findOrderBuilder{
		db:     db,
		req:    req,
		values: make([]any, 0),
	}
}

type findOrderEngineer struct {
	builder IFindOrderBuilder
}

func FundOrderEngineer(builder IFindOrderBuilder) *findOrderEngineer {
	return &findOrderEngineer{builder: builder}
}

func (b *findOrderBuilder) initQuery() {

	b.query += `
	SELECT 
		array_to_json(array_agg("at"))
	FROM (
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
		`
}

func (b *findOrderBuilder) initCountQuery() {

	b.query += `
	SELECT
	FROM "orders "o"
	WHERE 1 = 1`

}

func (b *findOrderBuilder) buildWhereSearch() {

	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		query := fmt.Sprintf(`
		AND (
			LOWER("o"."user_id") LIKE $%d OR
			LOWER("o"."address") LIKE $%d OR
			LOWER("o"."contact") LIKE $%d
		)`,
			b.lastIndex+1,
			b.lastIndex+2,
			b.lastIndex+3,
		)

		temp := b.getQuery()
		temp += query
		b.setQuery(query)

		b.lastIndex = len(b.values)
	}

}

func (b *findOrderBuilder) buildWhereStatus() {
	if b.req.Status != "" {
		b.values = append(
			b.values,
			strings.ToLower(b.req.Search),
		)

		query := fmt.Sprintf(`
		AND "o"."status" = $%d `,
			b.lastIndex+1,
		)

		temp := b.getQuery()
		temp += query
		b.setQuery(query)

		b.lastIndex = len(b.values)
	}
}

func (b *findOrderBuilder) buildWhereDate() {

	if b.req.StartDate != "" && b.req.EndDate != "" {
		b.values = append(
			b.values,
			b.req.StartDate,
			b.req.EndDate,
		)

		query := fmt.Sprintf(`
		AND "o"."created_date" BETWEEN ($%d)::DATE AND ($%d)::DATE + 1`,
			b.lastIndex+1,
			b.lastIndex+2,
		)

		temp := b.getQuery()
		temp += query
		b.setQuery(query)

		b.lastIndex = len(b.values)
	}

}

func (b *findOrderBuilder) builSort() {}

func (b *findOrderBuilder) buildPaginate() {}

func (b *findOrderBuilder) FinalQuery() {}

func (b *findOrderBuilder) closeQuery() {

	b.query += `
	) AS "at"`

}

func (b *findOrderBuilder) getQuery() string { return b.query }

func (b *findOrderBuilder) setQuery(query string) { b.query = query }

func (b *findOrderBuilder) getValues() []any { return b.values }

func (b *findOrderBuilder) setValues(data []any) { b.values = data }

func (b *findOrderBuilder) setLastIndex(n int) { b.lastIndex = n }

func (b *findOrderBuilder) getDb() *sqlx.DB { return b.db }

func (b *findOrderBuilder) reset() {

	b.query = ""
	b.values = make([]any, 0)
	b.lastIndex = 0
}

func (en *findOrderEngineer) FindOrder() []*orders.Order {
	return nil
}

func (en *findOrderEngineer) CountOrder() int {
	return 0
}
