package filesHandlers

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/guatom999/Ecommerce-Go/config"
	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/files"
	"github.com/guatom999/Ecommerce-Go/modules/files/filesUseCases"
	"github.com/guatom999/Ecommerce-Go/pkg/utils"
)

type fileHandlerErrCode string

const (
	uploadFileErr fileHandlerErrCode = "files-001"
)

type IFilesHandler interface {
	UploadFiles(c *fiber.Ctx) error
}

type filesHandler struct {
	cfg         config.IConfig
	fileUseCase filesUseCases.IFilesUseCases
}

func FilesHandler(cfg config.IConfig, fileUseCase filesUseCases.IFilesUseCases) IFilesHandler {
	return &filesHandler{
		cfg:         cfg,
		fileUseCase: fileUseCase,
	}
}

func (h *filesHandler) UploadFiles(c *fiber.Ctx) error {

	req := make([]*files.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadFileErr),
			err.Error(),
		).Res()
	}

	filesReq := form.File["files"]
	destination := c.FormValue("destination")

	// file extension validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}

	for _, file := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadFileErr),
				"files extension is not acceptable",
			).Res()
		}

		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadFileErr),
				// "file to large",
				fmt.Sprintf("file size must less than %d MiB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		filename := utils.RandFileName(ext)

		req = append(req, &files.FileReq{
			File:        file,
			Destination: destination + "/" + filename,
			Extension:   ext,
			FileName:    filename,
		})
	}

	res, err := h.fileUseCase.UploadToGCP(req)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(uploadFileErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		res,
	).Res()
}
