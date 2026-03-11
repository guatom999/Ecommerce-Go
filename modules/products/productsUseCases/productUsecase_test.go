package productsUseCases

import (
	"errors"
	"testing"

	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/products"
	"github.com/guatom999/Ecommerce-Go/modules/products/productsRepositories"
	"github.com/stretchr/testify/assert"
)

// helper: สร้าง useCase พร้อม mock repo
func newUseCaseWithMock() (IProductsUseCase, *productsRepositories.MockProductRepository) {
	mockRepo := new(productsRepositories.MockProductRepository)
	uc := ProductsUseCase(mockRepo)
	return uc, mockRepo
}

// --- FindOneProduct ---
// แยก 2 cases: repo สำเร็จ vs repo คืน error

func TestFindOneProduct_Success(t *testing.T) {
	uc, mockRepo := newUseCaseWithMock()

	expected := &products.Product{Id: "p1", Title: "Test Product"}
	mockRepo.On("FindOneProduct", "p1").Return(expected, nil)

	result, err := uc.FindOneProduct("p1")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestFindOneProduct_RepoError(t *testing.T) {
	uc, mockRepo := newUseCaseWithMock()

	mockRepo.On("FindOneProduct", "p1").Return((*products.Product)(nil), errors.New("not found"))

	result, err := uc.FindOneProduct("p1")

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// --- FindProduct ---
// ไม่มี error return → test เดียว: ตรวจว่า pagination คิดถูก

func TestFindProduct_Success(t *testing.T) {
	uc, mockRepo := newUseCaseWithMock()

	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{Page: 1, Limit: 10},
		SortReq:       &entities.SortReq{},
	}
	mockProducts := []*products.Product{
		{Id: "p1", Title: "A"},
		{Id: "p2", Title: "B"},
	}
	mockRepo.On("FindProduct", req).Return(mockProducts, 25)

	result := uc.FindProduct(req)

	assert.Equal(t, mockProducts, result.Data)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Limit)
	assert.Equal(t, 3, result.TotalPage) // ceil(25/10) = 3
	assert.Equal(t, 25, result.TotalItem)
	mockRepo.AssertExpectations(t)
}

// --- AddProduct ---
// แยก 2 cases: repo สำเร็จ vs repo คืน error

func TestAddProduct_Success(t *testing.T) {
	uc, mockRepo := newUseCaseWithMock()

	req := &products.Product{Title: "New Product", Price: 99.99}
	expected := &products.Product{Id: "p1", Title: "New Product", Price: 99.99}
	mockRepo.On("InsertProduct", req).Return(expected, nil)

	result, err := uc.AddProduct(req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestAddProduct_RepoError(t *testing.T) {
	uc, mockRepo := newUseCaseWithMock()

	req := &products.Product{Title: "New Product"}
	mockRepo.On("InsertProduct", req).Return((*products.Product)(nil), errors.New("insert failed"))

	result, err := uc.AddProduct(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// --- UpdateProduct ---
// แยก 2 cases: repo สำเร็จ vs repo คืน error

func TestUpdateProduct_Success(t *testing.T) {
	uc, mockRepo := newUseCaseWithMock()

	req := &products.Product{Id: "p1", Title: "Updated"}
	expected := &products.Product{Id: "p1", Title: "Updated"}
	mockRepo.On("UpdateProduct", req).Return(expected, nil)

	result, err := uc.UpdateProduct(req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_RepoError(t *testing.T) {
	uc, mockRepo := newUseCaseWithMock()

	req := &products.Product{Id: "p1", Title: "Updated"}
	mockRepo.On("UpdateProduct", req).Return((*products.Product)(nil), errors.New("update failed"))

	result, err := uc.UpdateProduct(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// --- DeleteProduct ---
// แยก 2 cases: repo สำเร็จ vs repo คืน error

func TestDeleteProduct_Success(t *testing.T) {
	uc, mockRepo := newUseCaseWithMock()

	mockRepo.On("DeleteProduct", "p1").Return(nil)

	err := uc.DeleteProduct("p1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct_RepoError(t *testing.T) {
	uc, mockRepo := newUseCaseWithMock()

	mockRepo.On("DeleteProduct", "p1").Return(errors.New("delete failed"))

	err := uc.DeleteProduct("p1")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
