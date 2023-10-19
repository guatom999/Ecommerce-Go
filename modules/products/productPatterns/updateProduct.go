package productPatterns

import (
	"context"
	"fmt"

	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/files"
	"github.com/guatom999/Ecommerce-Go/modules/files/filesUseCases"
	"github.com/guatom999/Ecommerce-Go/modules/products"
	"github.com/jmoiron/sqlx"
)

type IUpdateProductBuilder interface {
	initTransaction() error
	initQuery()
	updateTitleQuery()
	updateDescriptionQuery()
	updatePriceQuery()
	updateCategory() error
	insertImages() error
	getOldImages() []*entities.Image
	deleteOldImages() error
	closeQuery()
	updateProduct() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
	// show all query
	printQuery() string
}

type updateproductBuilder struct {
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *products.Product
	filesUsecases  filesUseCases.IFilesUseCases
	query          string
	queryField     []string
	lastStackIndex int
	values         []any
}

func UpdateProductBuilder(db *sqlx.DB, req *products.Product, filesUsecases filesUseCases.IFilesUseCases) IUpdateProductBuilder {
	return &updateproductBuilder{
		db:            db,
		req:           req,
		filesUsecases: filesUsecases,
		queryField:    make([]string, 0),
		values:        make([]any, 0),
		// lastStackIndex: 0,
	}
}

func (b *updateproductBuilder) initTransaction() error {

	tx, err := b.db.BeginTxx(context.Background(), nil)

	if err != nil {
		return err
	}

	b.tx = tx

	return nil

}
func (b *updateproductBuilder) initQuery() {

	b.query += `
	UPDATE "products" SET`

}
func (b *updateproductBuilder) updateTitleQuery() {

	if b.req.Title != "" {
		b.values = append(b.values, b.req.Title)
		b.lastStackIndex = len(b.values)

		b.queryField = append(b.queryField, fmt.Sprintf(`
		"title" = $%d`, b.lastStackIndex))

	}

}
func (b *updateproductBuilder) updateDescriptionQuery() {

	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		b.lastStackIndex = len(b.values)

		b.queryField = append(b.queryField, fmt.Sprintf(`
		"description" = $%d`, b.lastStackIndex))

	}

}
func (b *updateproductBuilder) updatePriceQuery() {

	if b.req.Price > 0 {
		b.values = append(b.values, b.req.Price)
		b.lastStackIndex = len(b.values)

		b.queryField = append(b.queryField, fmt.Sprintf(`
		"price" = $%d`, b.lastStackIndex))
	}

}
func (b *updateproductBuilder) updateCategory() error {

	if b.req.Category == nil {
		return nil
	}

	if b.req.Category.Id == 0 {
		return nil
	}

	query := `
	UPDATE "products_categories" SET 
		"category_id" = $1 
	WHERE "product_id" = $2;`

	if _, err := b.tx.ExecContext(context.Background(),
		query,
		b.req.Category.Id,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update product_categories failed: %v", err)
	}

	return nil

}
func (b *updateproductBuilder) insertImages() error {
	query := `
	INSERT INTO "images" (
		"filename",
		"url",
		"product_id"
	)
	VALUES`

	valueStack := make([]any, 0)
	var index int
	for i := range b.req.Image {
		valueStack = append(valueStack,
			b.req.Image[i].FileName,
			b.req.Image[i].Url,
			b.req.Id,
		)

		if i != len(b.req.Image)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d);`, index+1, index+2, index+3)
		}
		index += 3
	}

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}
	return nil
}

func (b *updateproductBuilder) getOldImages() []*entities.Image {

	query := `
	SELECT 
		"id",
		"filename",
		"url"
	FROM "images"
	WHERE "product_id" = $1;`

	images := make([]*entities.Image, 0)
	if err := b.db.Select(&images, query, b.req.Id); err != nil {
		fmt.Printf("get old image failed: %v", err)
		return make([]*entities.Image, 0)
	}

	return images
}
func (b *updateproductBuilder) deleteOldImages() error {

	query := `
	DELETE FROM "images" 
	WHERE "product_id" = $1;`

	images := b.getOldImages()

	if len(images) > 0 {
		deleteFileReq := make([]*files.DeleteFileReq, 0)
		for _, image := range images {
			deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
				Destination: fmt.Sprintf("images/products/%s", image.FileName),
			})
		}
		b.filesUsecases.DeleteFileOnGCP(deleteFileReq)

	}
	if _, err := b.tx.ExecContext(context.Background(), query, b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("delete imaged failed:%v", err)
	}

	return nil

}
func (b *updateproductBuilder) closeQuery() {

	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	// b.queryField = append(b.queryField, fmt.Sprintf(`"price" = $%d`, b.lastStackIndex))
	b.query += fmt.Sprintf(` 
	WHERE "id" =  $%d`, b.lastStackIndex)

}

func (b *updateproductBuilder) printQuery() string {

	if b.query == "" {
		return fmt.Sprint("no query right now")
	}

	return fmt.Sprintf("query string is :%v", b.query)

}

func (b *updateproductBuilder) updateProduct() error {

	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update product failed:%v", err)
	}

	return nil

}
func (b *updateproductBuilder) getQueryFields() []string { return b.queryField }
func (b *updateproductBuilder) getValues() []any         { return b.values }
func (b *updateproductBuilder) getQuery() string         { return b.query }
func (b *updateproductBuilder) setQuery(query string)    { b.query = query }
func (b *updateproductBuilder) getImagesLen() int        { return len(b.req.Image) }
func (b *updateproductBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return err
	}
	return nil
}

type updateProductEngineer struct {
	builder IUpdateProductBuilder
}

func UpdateProductEngineer(builder IUpdateProductBuilder) *updateProductEngineer {
	return &updateProductEngineer{builder: builder}
}

func (en *updateProductEngineer) sumQueryFeilds() {

	en.builder.updateTitleQuery()
	en.builder.updateDescriptionQuery()
	en.builder.updatePriceQuery()

	fields := en.builder.getQueryFields()

	for i := range fields {
		query := en.builder.getQuery()
		if i != len(fields)-1 {
			en.builder.setQuery(query + fields[i] + ",")
		} else {
			en.builder.setQuery(query + fields[i])
		}
	}
	// en.builder.updateCategoryQuery()

}

func (en *updateProductEngineer) UpdateProduct() error {

	en.builder.initTransaction()

	en.builder.initQuery()
	en.sumQueryFeilds()
	en.builder.closeQuery()

	fmt.Println(en.builder.getQuery())

	if err := en.builder.updateProduct(); err != nil {
		return err
	}

	if err := en.builder.updateCategory(); err != nil {
		return err
	}

	if en.builder.getImagesLen() > 0 {
		if err := en.builder.deleteOldImages(); err != nil {
			return err
		}
		if err := en.builder.insertImages(); err != nil {
			return err
		}
	}

	if err := en.builder.commit(); err != nil {
		return err
	}

	return nil
}
