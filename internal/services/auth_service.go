package service

import (
	"errors"
	"log"
	"time"

	"github.com/KsenoTech/merch-store/internal/models"
	"github.com/KsenoTech/merch-store/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	secretKey string
	UserRepo  *repository.UserRepository
}

func NewAuthService(secretKey string, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{secretKey: secretKey, UserRepo: userRepo}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *AuthService) ComparePasswords(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func (s *AuthService) GenerateToken(userID uint) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *AuthService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *AuthService) AuthenticateOrCreateUser(username, password string) (string, error) {
	log.Printf("Authenticating or creating user: %s", username)

	// Проверяем, существует ли пользователь
	user, err := s.UserRepo.GetUserByUsername(username)
	if err != nil {
		log.Printf("User %s not found, creating...", username)

		// Хешируем пароль
		hashedPassword, err := s.HashPassword(password)
		if err != nil {
			return "", errors.New("failed to hash password")
		}

		// Создаем нового пользователя
		newUser := models.User{
			Username: username,
			Password: hashedPassword,
			Coins:    1000, // Начальный баланс
		}
		if err := s.UserRepo.CreateUser(&newUser); err != nil {
			return "", errors.New("failed to create user")
		}

		// Генерируем токен для нового пользователя
		token, err := s.GenerateToken(newUser.ID)
		if err != nil {
			return "", errors.New("failed to generate token")
		}
		log.Printf("User %s created successfully with ID: %d", username, newUser.ID)
		return token, nil
	}

	// Пользователь существует, проверяем пароль
	if !s.ComparePasswords(user.Password, password) { // Используем результат ComparePasswords напрямую
		return "", errors.New("invalid password")
	}

	// Генерируем токен для существующего пользователя
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	return token, nil
}
