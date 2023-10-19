package productHandlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/appinfo"
	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/files"
	"github.com/guatom999/Ecommerce-Go/modules/files/filesUseCases"
	"github.com/guatom999/Ecommerce-Go/modules/products"
	"github.com/guatom999/Ecommerce-Go/modules/products/productsUseCases"
)

type productErrCode string

const (
	findOneProductErr productErrCode = "product-001"
	findProductErr    productErrCode = "product-002"
	insertProductErr  productErrCode = "product-003"
	deleteProductErr  productErrCode = "product-004"
	updateProductErr  productErrCode = "product-005"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindProduct(c *fiber.Ctx) error
	AddProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type productsHandler struct {
	cfg             config.IConfig
	productsUseCase productsUseCases.IProductsUseCase
	filesUseCases   filesUseCases.IFilesUseCases
}

func ProductsHandler(cfg config.IConfig, productsUseCase productsUseCases.IProductsUseCase, filesUseCases filesUseCases.IFilesUseCases) IProductsHandler {
	return &productsHandler{
		cfg:             cfg,
		productsUseCase: productsUseCase,
		filesUseCases:   filesUseCases,
	}
}

func (h *productsHandler) FindOneProduct(c *fiber.Ctx) error {

	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productsUseCase.FindOneProduct(productId)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findOneProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		product,
	).Res()
}

func (h *productsHandler) FindProduct(c *fiber.Ctx) error {

	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findProductErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}

	if req.Sort == "" {
		req.Sort = "ASC"
	}

	product := h.productsUseCase.FindProduct(req)

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		product,
	).Res()

}

func (h *productsHandler) AddProduct(c *fiber.Ctx) error {

	req := &products.Product{
		Category: &appinfo.Category{},
		Image:    make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}
	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertProductErr),
			"category id is invalid",
		).Res()
	}

	product, err := h.productsUseCase.AddProduct(req)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadGateway.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		product,
	).Res()
}

func (h *productsHandler) UpdateProduct(c *fiber.Ctx) error {

	productId := strings.Trim(c.Params("product_id"), " ")
	req := &products.Product{
		Category: &appinfo.Category{},
		Image:    make([]*entities.Image, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateProductErr),
			err.Error(),
		).Res()
	}

	req.Id = productId

	product, err := h.productsUseCase.UpdateProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		product,
	).Res()
}

func (h *productsHandler) DeleteProduct(c *fiber.Ctx) error {

	productId := strings.Trim(c.Params("product_id"), " ")
	product, err := h.productsUseCase.FindOneProduct(productId)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, file := range product.Image {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("images/test/%s", file.FileName),
		})
	}

	if err := h.filesUseCases.DeleteFileOnGCP(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	if err := h.productsUseCase.DeleteProduct(productId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		string("delete success"),
	).Res()
}
