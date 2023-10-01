package appinfohandlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/appinfo"
	appinfousecases "github.com/guatom999/Ecommerce-Go/modules/appinfo/appinfoUseCases"
	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/pkg/authen"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	AddCategory(c *fiber.Ctx) error
	RemoveCategory(c *fiber.Ctx) error
}

type appinfoHandler struct {
	config         config.IConfig
	appinfoUsecase appinfousecases.IAppinfoUsecase
}

type appinfoHandlerErrorCode string

const (
	generateApiKeyErr appinfoHandlerErrorCode = "appinfo-001"
	findCategoryErr   appinfoHandlerErrorCode = "appinfo-002"
	addCategoryErr    appinfoHandlerErrorCode = "appinfo-003"
	removeCategoryErr appinfoHandlerErrorCode = "appinfo-004"
)

func AppinfoHandler(config config.IConfig, appinfoUsecase appinfousecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		config:         config,
		appinfoUsecase: appinfoUsecase,
	}
}

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apikey, err := authen.NewAuth(
		authen.ApiKey,
		h.config.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(generateApiKeyErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,

		&struct {
			Key string `json:"key"`
		}{
			Key: apikey.SignToken(),
		},
	).Res()
}

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	category, err := h.appinfoUsecase.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		category,
	).Res()
}

func (h *appinfoHandler) AddCategory(c *fiber.Ctx) error {

	req := make([]*appinfo.Category, 0)
	// req := new([]*appinfo.Category)
	if err := c.BodyParser(&req); err != nil {
		fmt.Println("Error Here1")
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}

	if len(req) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			"body req is empthy",
		).Res()
	}

	err := h.appinfoUsecase.InsertCategory(req)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		nil,
	).Res()
}

func (h *appinfoHandler) RemoveCategory(c *fiber.Ctx) error {

	categoryId := strings.Trim(c.Params("category_id"), " ")
	categoryIdInt, err := strconv.Atoi(categoryId)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(removeCategoryErr),
			"category id type is invalid",
		).Res()
	}

	if categoryIdInt <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(removeCategoryErr),
			"id must more than 0",
		).Res()
	}

	if err := h.appinfoUsecase.DeleteCategory(categoryIdInt); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(removeCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			CategoryId int `json:"category_id" `
		}{
			CategoryId: categoryIdInt,
		},
	).Res()
}
