package productPatterns

import (
	"context"
	"fmt"
	"time"

	"github.com/guatom999/Ecommerce-Go/modules/products"
	"github.com/jmoiron/sqlx"
)

type IInsertProductBuilder interface {
	initTransaction() error
	insertProduct() error
	insertCategory() error
	insertAttacthment() error
	commit() error
	getProductId() string
}

type insertProductBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *products.Product
}

// constructor
func InsertProductBuilder(db *sqlx.DB, req *products.Product) IInsertProductBuilder {
	return &insertProductBuilder{
		db:  db,
		req: req,
	}
}

// |----------------------------------------------------------------------------------------------------------|
// Engineer PaDB
type insertProductEngineer struct {
	builder IInsertProductBuilder
}

func (b *insertProductBuilder) initTransaction() error {

	tx, err := b.db.BeginTxx(context.Background(), nil)

	if err != nil {
		return err
	}

	b.tx = tx

	return nil
}
func (b *insertProductBuilder) insertProduct() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)

	defer cancel()

	query := `
	INSERT INTO "products" ( 
		"title" , 
		"description" , 
		"price" 
	)
	VALUES ( 
		$1 , 
		$2 , 
		$3
	)
	RETURNING "id";
	`

	if err := b.tx.QueryRowContext(
		ctx,
		query,
		b.req.Title,
		b.req.Description,
		b.req.Price,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert product failed: %v", err)
	}

	return nil
}
func (b *insertProductBuilder) insertCategory() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "products_categories" ( 
		"product_id",
		"category_id",
	) VALUES ($1 , $2);
	`

	if _, err := b.tx.ExecContext(ctx, query, b.req.Id, b.req.Category.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert product category failed: %v", err)
	}

	return nil
}
func (b *insertProductBuilder) insertAttacthment() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "images" ( 
		"filename",
		"url",
		"product_id"
	)
	VAULES
	`

	valuesStack := make([]any, 0)
	var index int
	for i := range b.req.Image {
		valuesStack = append(
			valuesStack,
			b.req.Image[i].FileName,
			b.req.Image[i].Url,
			b.req.Id,
		)

		if i != len(b.req.Image)-1 {
			query += fmt.Sprintf(`($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`($%d, $%d, $%d);`, index+1, index+2, index+3)
		}

		index += 3
	}

	if _, err := b.tx.ExecContext(
		ctx,
		query,
		valuesStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert image failed: %v", err)
	}

	return nil
}
func (b *insertProductBuilder) commit() error {
	return nil
}

func (b *insertProductBuilder) getProductId() string {
	return ""
}

func InsertProductEngineer(builder IInsertProductBuilder) *insertProductEngineer {
	return &insertProductEngineer{builder: builder}
}

func (en *insertProductEngineer) InsertProduct() (string, error) {
	return "", nil
}
