package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cookit/backend/models"
	"github.com/cookit/backend/repositories"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_COST = 8
)

type UserService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUserById(ctx *gin.Context, id int) (*models.User, error) {
	return s.userRepository.GetUserByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx *gin.Context, email string) (*models.User, error) {
	return s.userRepository.GetUserByEmail(ctx, email)
}

func (s *UserService) CreateUser(ctx *gin.Context, user *models.User) (*models.User, error) {
	repeatedUser, err := s.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if repeatedUser.Id != "" {
		return nil, errors.New("User already exists`")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), HASH_COST)
	if err != nil {
		return nil, err
	}

	err = s.userRepository.CreateUser(ctx, user.Email, string(password))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(ctx *gin.Context, email, password string) (*models.LoginResponse, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user.Id == "" {
		return nil, errors.New("Invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println(err)
		return nil, errors.New("Invalid credentials")
	}
	claims := models.AppClaims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}
	return &models.LoginResponse{
		Token: tokenString,
	}, nil
}
