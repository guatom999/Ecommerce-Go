package middlewaresRepositories

import "github.com/jmoiron/sqlx"

type IMiddlewareRepository interface {
	FindAccessToken(userId string, accessToken string) bool
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
