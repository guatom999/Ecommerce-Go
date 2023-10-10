package productsRepositories

import (
	"encoding/json"
	"fmt"

	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/files/filesUseCases"
	"github.com/guatom999/Ecommerce-Go/modules/products"
	"github.com/guatom999/Ecommerce-Go/modules/products/productPatterns"
	"github.com/jmoiron/sqlx"
)

type IProductRepository interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProduct(req *products.ProductFilter) ([]*products.Product, int)
}

type productRepository struct {
	db          *sqlx.DB
	cfg         config.IConfig
	fileUsecase filesUseCases.IFilesUseCases
}

func ProductRepository(db *sqlx.DB, cfg config.IConfig, fileUsecase filesUseCases.IFilesUseCases) IProductRepository {
	return &productRepository{
		db:          db,
		cfg:         cfg,
		fileUsecase: fileUsecase,
	}
}

func (r *productRepository) FindOneProduct(productId string) (*products.Product, error) {
	query :=
		`
		SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"p"."id",
			"p"."title",
			"p"."description",
			"p"."price",
			(
				SELECT
					to_jsonb("ct")
				FROM (
					SELECT
						"c"."id",
						"c"."title"
					FROM "categories" "c"
						LEFT JOIN "products_categories" "pc" ON "pc"."category_id" = "c"."id"
					WHERE "pc"."product_id" = "p"."id"
				) AS "ct"
			) AS "category",
			"p"."created_at",
			"p"."updated_at",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "images" "i"
					WHERE "i"."product_id" = "p"."id"
				) AS "it"
			) AS "images"
		FROM "products" "p"
		WHERE "p"."id" = $1
		LIMIT 1
	) AS "t";
	`

	productByte := make([]byte, 0)
	product := &products.Product{
		Image: make([]*entities.Image, 0),
	}

	if err := r.db.Get(&productByte, query, productId); err != nil {
		return nil, fmt.Errorf("get product failed: %v", err)
	}

	if err := json.Unmarshal(productByte, &product); err != nil {
		return nil, fmt.Errorf("unmarshal product failed: %v", err)
	}

	return product, nil
}

func (r *productRepository) FindProduct(req *products.ProductFilter) ([]*products.Product, int) {
	builder := productPatterns.FindProductBuilder(r.db, req)
	engineer := productPatterns.FindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()

	// fmt.Printf("result is :%v", result)

	return result, count

}
