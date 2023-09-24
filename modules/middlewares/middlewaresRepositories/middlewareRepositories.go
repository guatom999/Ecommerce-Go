package middlewaresRepositories

import (
	"fmt"

	"github.com/guatom999/Ecommerce-Go/modules/middlewares"
	"github.com/jmoiron/sqlx"
)

type IMiddlewareRepository interface {
	FindAccessToken(userId string, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewareRepository struct {
	db *sqlx.DB
}

func MiddlewareRepository(db *sqlx.DB) IMiddlewareRepository {
	return &middlewareRepository{db: db}
}

func (m *middlewareRepository) FindAccessToken(userId string, accessToken string) bool {

	query := `
	SELECT 
		(CASE WHEN COUNT(*) = 1 THEN TRUE ELSE FALSE END) 
	FROM "oauth"
	WHERE "user_id" = $1 
	AND "access_token" = $2;
	`

	var check bool
	if err := m.db.Get(&check, query, userId, accessToken); err != nil {
		return false
	}

	return true
}

func (m *middlewareRepository) FindRole() ([]*middlewares.Role, error) {
	query := `SELECT "id" , "title" FROM "roles" ORDER BY "id" DESC `

	roles := make([]*middlewares.Role, 0)
	if err := m.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("roles are empthy")
	}

	return roles, nil
}
