package servers

import (
	"github.com/gofiber/fiber/v2"
	appinfohandlers "github.com/guatom999/Ecommerce-Go/modules/appinfo/appinfoHandlers"
	appinforepositories "github.com/guatom999/Ecommerce-Go/modules/appinfo/appinfoRepositories"
	appinfousecases "github.com/guatom999/Ecommerce-Go/modules/appinfo/appinfoUseCases"
	"github.com/guatom999/Ecommerce-Go/modules/files/filesHandlers"
	"github.com/guatom999/Ecommerce-Go/modules/files/filesUseCases"
	"github.com/guatom999/Ecommerce-Go/modules/middlewares/middlewaresHandlers"
	"github.com/guatom999/Ecommerce-Go/modules/middlewares/middlewaresRepositories"
	"github.com/guatom999/Ecommerce-Go/modules/middlewares/middlewaresUsecases"
	"github.com/guatom999/Ecommerce-Go/modules/monitor/monitorHandlers"
	"github.com/guatom999/Ecommerce-Go/modules/orders/ordersHandlers"
	"github.com/guatom999/Ecommerce-Go/modules/orders/ordersRepositories"
	"github.com/guatom999/Ecommerce-Go/modules/orders/ordersUseCases"
	"github.com/guatom999/Ecommerce-Go/modules/products/productHandlers"
	"github.com/guatom999/Ecommerce-Go/modules/products/productsRepositories"
	"github.com/guatom999/Ecommerce-Go/modules/products/productsUseCases"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersHandlers"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersRepositories"
	"github.com/guatom999/Ecommerce-Go/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppInfoModule()
	FileTransferModule()
	ProductsModule()
	OrderModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
	mid    middlewaresHandlers.IMiddlewareHandler
}

// Constructor
func NewModule(router fiber.Router, server *server, mid middlewaresHandlers.IMiddlewareHandler) IModuleFactory {
	return &moduleFactory{
		router: router,
		server: server,
		mid:    mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewareHandler {
	repository := middlewaresRepositories.MiddlewareRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	handler := middlewaresHandlers.MiddlewareHandler(s.cfg, usecase)

	return handler
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.server.db)
	usecase := usersUsecases.UsersUsecase(m.server.cfg, repository)
	handler := usersHandlers.UsersHandler(m.server.cfg, usecase)

	router := m.router.Group("/users")

	router.Post("/signup", m.mid.ApiKeyCheck(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyCheck(), handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyCheck(), handler.RefeshPassport)
	router.Post("/signout", m.mid.ApiKeyCheck(), handler.SignOut)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignUpAdmin)

	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)

}

func (m *moduleFactory) AppInfoModule() {
	repository := appinforepositories.AppinfoRepository(m.server.db)
	usecase := appinfousecases.AppinfoUseCase(repository)
	handler := appinfohandlers.AppinfoHandler(m.server.cfg, usecase)

	router := m.router.Group("/appinfo")

	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)
	router.Get("/categories", m.mid.ApiKeyCheck(), handler.FindCategory)

	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)

	router.Delete("/:category_id/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCategory)

}

func (m *moduleFactory) FileTransferModule() {
	usecase := filesUseCases.FilesUseCase(m.server.cfg)
	handler := filesHandlers.FilesHandler(m.server.cfg, usecase)

	router := m.router.Group("/files")

	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UploadFiles)
	router.Patch("/delete", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteFile)
	// _ = router
	// _ = handler
}

func (m *moduleFactory) ProductsModule() {
	fileUsecase := filesUseCases.FilesUseCase(m.server.cfg)

	productsRepository := productsRepositories.ProductRepository(m.server.db, m.server.cfg, fileUsecase)
	productsUsecase := productsUseCases.ProductsUseCase(productsRepository)
	productsHandler := productHandlers.ProductsHandler(m.server.cfg, productsUsecase, fileUsecase)

	router := m.router.Group("/products")

	router.Post("/", m.mid.JwtAuth(), m.mid.Authorize(2), productsHandler.AddProduct)
	router.Delete("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), productsHandler.DeleteProduct)

	router.Patch("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), productsHandler.UpdateProduct)

	router.Get("/", m.mid.ApiKeyCheck(), productsHandler.FindProduct)
	router.Get("/:product_id", m.mid.ApiKeyCheck(), productsHandler.FindOneProduct)

	// router.Delete("/:product_id", m.mid.ApiKeyCheck(), productsHandler.DeleteProduct)

}

func (m *moduleFactory) OrderModule() {

	fileUsecase := filesUseCases.FilesUseCase(m.server.cfg)
	productsRepository := productsRepositories.ProductRepository(m.server.db, m.server.cfg, fileUsecase)

	ordersRepository := ordersRepositories.OrderRepository(m.server.db)
	ordersUsecase := ordersUseCases.OrderUseCase(ordersRepository, productsRepository)
	ordersHandler := ordersHandlers.OrderHandler(m.server.cfg, ordersUsecase)

	router := m.router.Group("/orders")

	router.Get("/:order_id", m.mid.JwtAuth(), ordersHandler.FindOneOrder)
	router.Get("/", m.mid.JwtAuth(), m.mid.ParamsCheck(), ordersHandler.FindOrder)

	router.Post("/", m.mid.JwtAuth(), ordersHandler.InsertOrder)

}
