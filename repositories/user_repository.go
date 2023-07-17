package repositories

import (
	"github.com/cookit/backend/models"
	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	GetUserByID(ctx *gin.Context, userId int) (*models.User, error)
	GetUserByEmail(ctx *gin.Context, email string) (*models.User, error)
	CreateUser(ctx *gin.Context, email, password string) error
}
