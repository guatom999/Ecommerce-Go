package ordersUseCases

import (
	"errors"
	"testing"

	"github.com/guatom999/Ecommerce-Go/modules/entities"
	"github.com/guatom999/Ecommerce-Go/modules/orders"
	"github.com/guatom999/Ecommerce-Go/modules/orders/ordersRepositories"
	"github.com/guatom999/Ecommerce-Go/modules/products"
	"github.com/guatom999/Ecommerce-Go/modules/products/productsRepositories"
	"github.com/stretchr/testify/assert"
)

func newOrderUseCaseWithMock() (IOrderUseCase, *ordersRepositories.MockOrderRepository, *productsRepositories.MockProductRepository) {
	orderMockRepo := new(ordersRepositories.MockOrderRepository)
	productMockRepo := new(productsRepositories.MockProductRepository)
	uc := OrderUseCase(orderMockRepo, productMockRepo)
	return uc, orderMockRepo, productMockRepo
}

// helper สำหรับ test ที่ต้องการ productMockRepo ด้วย (InsertOrder)
// func newOrderUseCaseWithBothMocks() (IOrderUseCase, *ordersRepositories.MockOrderRepository, *productsRepositories.MockProductRepository) {
// 	orderMockRepo := new(ordersRepositories.MockOrderRepository)
// 	productMockRepo := new(productsRepositories.MockProductRepository)
// 	uc := OrderUseCase(orderMockRepo, productMockRepo)
// 	return uc, orderMockRepo, productMockRepo
// }

// --- FindOneOrder ---

func TestFindOneOrder_Success(t *testing.T) {
	uc, orderMockRepo, _ := newOrderUseCaseWithMock()

	expected := &orders.Order{Id: "order1", UserId: "user1"}
	orderMockRepo.On("FindOneOrder", "order1").Return(expected, nil)

	result, err := uc.FindOneOrder("order1")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	orderMockRepo.AssertExpectations(t)
}

func TestFindOneOrder_Failed(t *testing.T) {
	uc, orderMockRepo, _ := newOrderUseCaseWithMock()

	orderMockRepo.On("FindOneOrder", "order1").Return((*orders.Order)(nil), errors.New("order not found"))

	result, err := uc.FindOneOrder("order1")

	assert.Error(t, err)
	assert.Nil(t, result)
	orderMockRepo.AssertExpectations(t)
}

// --- FindOrder ---
// ไม่มี error return → test เดียว: ตรวจ pagination

func TestFindOrder_Success(t *testing.T) {
	uc, orderMockRepo, _ := newOrderUseCaseWithMock()

	req := &orders.OrderFilter{
		PaginationReq: &entities.PaginationReq{Page: 1, Limit: 5},
		SortReq:       &entities.SortReq{},
	}
	mockOrders := []*orders.Order{
		{Id: "order1", UserId: "user1"},
		{Id: "order2", UserId: "user2"},
	}
	orderMockRepo.On("FindOrder", req).Return(mockOrders, 12)

	result := uc.FindOrder(req)

	assert.Equal(t, mockOrders, result.Data)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 5, result.Limit)
	assert.Equal(t, 3, result.TotalPage) // ceil(12/5) = 3
	assert.Equal(t, 12, result.TotalItem)
	orderMockRepo.AssertExpectations(t)
}

// --- InsertOrder ---
// มี 3 จุด error:
//   1. product เป็น nil (ก่อน mock ถูกเรียก)
//   2. productRepo.FindOneProduct ล้มเหลว
//   3. orderRepo.InsertOrder ล้มเหลว
//   4. orderRepo.FindOneOrder (หลัง insert) ล้มเหลว

func TestInsertOrder_ProductIsNil(t *testing.T) {
	uc, _, _ := newOrderUseCaseWithMock()

	req := &orders.Order{
		Product: []*orders.ProductsOrder{
			{Product: nil, Quantity: 1}, // product เป็น nil
		},
	}

	result, err := uc.InsertOrder(req)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestInsertOrder_FindOneProductFailed(t *testing.T) {
	uc, _, productMockRepo := newOrderUseCaseWithMock()

	req := &orders.Order{
		Product: []*orders.ProductsOrder{
			{Product: &products.Product{Id: "p1"}, Quantity: 2},
		},
	}
	productMockRepo.On("FindOneProduct", "p1").Return((*products.Product)(nil), errors.New("product not found"))

	result, err := uc.InsertOrder(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	productMockRepo.AssertExpectations(t)
}

func TestInsertOrder_InsertOrderFailed(t *testing.T) {
	uc, orderMockRepo, productMockRepo := newOrderUseCaseWithMock()

	req := &orders.Order{
		Product: []*orders.ProductsOrder{
			{Product: &products.Product{Id: "p1", Price: 100}, Quantity: 2},
		},
	}
	productMockRepo.On("FindOneProduct", "p1").Return(&products.Product{Id: "p1", Price: 100}, nil)
	orderMockRepo.On("InsertOrder", req).Return("", errors.New("insert failed"))

	result, err := uc.InsertOrder(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	productMockRepo.AssertExpectations(t)
	orderMockRepo.AssertExpectations(t)
}

func TestInsertOrder_FindOneOrderAfterInsertFailed(t *testing.T) {
	uc, orderMockRepo, productMockRepo := newOrderUseCaseWithMock()

	req := &orders.Order{
		Product: []*orders.ProductsOrder{
			{Product: &products.Product{Id: "p1", Price: 100}, Quantity: 1},
		},
	}
	productMockRepo.On("FindOneProduct", "p1").Return(&products.Product{Id: "p1", Price: 100}, nil)
	orderMockRepo.On("InsertOrder", req).Return("order1", nil)
	orderMockRepo.On("FindOneOrder", "order1").Return((*orders.Order)(nil), errors.New("find failed"))

	result, err := uc.InsertOrder(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	productMockRepo.AssertExpectations(t)
	orderMockRepo.AssertExpectations(t)
}

func TestInsertOrder_Success(t *testing.T) {
	uc, orderMockRepo, productMockRepo := newOrderUseCaseWithMock()

	req := &orders.Order{
		Product: []*orders.ProductsOrder{
			{Product: &products.Product{Id: "p1", Price: 100}, Quantity: 2},
		},
	}
	fullProduct := &products.Product{Id: "p1", Price: 100, Title: "Item A"}
	expected := &orders.Order{Id: "order1", TotalPaid: 200}

	productMockRepo.On("FindOneProduct", "p1").Return(fullProduct, nil)
	orderMockRepo.On("InsertOrder", req).Return("order1", nil)
	orderMockRepo.On("FindOneOrder", "order1").Return(expected, nil)

	result, err := uc.InsertOrder(req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	productMockRepo.AssertExpectations(t)
	orderMockRepo.AssertExpectations(t)
}

// --- UpdateOrder ---
// มี 2 จุด error:
//   1. orderRepo.UpdateOrder ล้มเหลว
//   2. orderRepo.FindOneOrder (หลัง update) ล้มเหลว

func TestUpdateOrder_UpdateFailed(t *testing.T) {
	uc, orderMockRepo, _ := newOrderUseCaseWithMock()

	req := &orders.Order{Id: "order1", Status: "paid"}
	orderMockRepo.On("UpdateOrder", req).Return(errors.New("update failed"))

	result, err := uc.UpdateOrder(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	orderMockRepo.AssertExpectations(t)
}

func TestUpdateOrder_FindOneOrderAfterUpdateFailed(t *testing.T) {
	uc, orderMockRepo, _ := newOrderUseCaseWithMock()

	req := &orders.Order{Id: "order1", Status: "paid"}
	orderMockRepo.On("UpdateOrder", req).Return(nil)
	orderMockRepo.On("FindOneOrder", "order1").Return((*orders.Order)(nil), errors.New("find failed"))

	result, err := uc.UpdateOrder(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	orderMockRepo.AssertExpectations(t)
}

func TestUpdateOrder_Success(t *testing.T) {
	uc, orderMockRepo, _ := newOrderUseCaseWithMock()

	req := &orders.Order{Id: "order1", Status: "paid"}
	expected := &orders.Order{Id: "order1", Status: "paid"}
	orderMockRepo.On("UpdateOrder", req).Return(nil)
	orderMockRepo.On("FindOneOrder", "order1").Return(expected, nil)

	result, err := uc.UpdateOrder(req)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	orderMockRepo.AssertExpectations(t)
}
