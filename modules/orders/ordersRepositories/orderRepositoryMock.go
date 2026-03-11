package ordersRepositories

import (
	"github.com/guatom999/Ecommerce-Go/modules/orders"
	"github.com/stretchr/testify/mock"
)

type (
	MockOrderRepository struct {
		mock.Mock
	}
)

func (m *MockOrderRepository) FindOneOrder(orderId string) (*orders.Order, error) {

	args := m.Called(orderId)
	return args.Get(0).(*orders.Order), args.Error(1)
}
func (m *MockOrderRepository) FindOrder(req *orders.OrderFilter) ([]*orders.Order, int) {
	args := m.Called(req)
	return args.Get(0).([]*orders.Order), args.Int(1)
}
func (m *MockOrderRepository) InsertOrder(req *orders.Order) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)

}
func (m *MockOrderRepository) UpdateOrder(req *orders.Order) error {
	args := m.Called(req)
	return args.Error(0)
}
