package productHandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/files/filesUseCases"
	"github.com/guatom999/Ecommerce-Go/modules/products/productsUseCases"
)

type productErrCode string

const (
	findOneProductErr productErrCode = "product-001"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
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
