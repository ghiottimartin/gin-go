package repositories

import (
	"database/sql"
	"log"

	"github.com/cookit/backend/models"
	"github.com/gin-gonic/gin"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db,
	}
}

func (r *PostgresUserRepository) GetUserByID(ctx *gin.Context, userId int) (*models.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) CreateUser(ctx *gin.Context, email, password string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (email, password) VALUES ($1, $2)", email, password)
	return err
}

func (r *PostgresUserRepository) GetUserByEmail(ctx *gin.Context, email string) (*models.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var user = models.User{}
	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}
