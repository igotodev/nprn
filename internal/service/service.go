package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/golang-jwt/jwt"
	"nprn/internal/customerr"
	"nprn/internal/entity/sale/salestorage"
	"nprn/internal/entity/user/usermodel"
	"nprn/internal/entity/user/userstorage"
	"nprn/pkg/logging"
	"time"
)

const (
	salt      = "4hsd83jd7fsd2"
	tokenTime = 12 * time.Hour
	signKey   = "dkr3!#mc349x#s3&74f12d"
)

type Service struct {
	userStorage userstorage.UserStorage
	saleStorage salestorage.SaleStorage
	logger      *logging.Logger
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

func NewService(userStorage userstorage.UserStorage, saleStorage salestorage.SaleStorage, logger *logging.Logger) *Service {
	return &Service{
		userStorage: userStorage,
		saleStorage: saleStorage,
		logger:      logger,
	}
}

func (s *Service) SignUp(ctx context.Context, user usermodel.UserInternal) (string, error) {
	passwordHash := s.generatePasswordHash(user.PasswordHash)
	user.PasswordHash = passwordHash

	objID, err := s.userStorage.Create(ctx, user)
	if err != nil {
		s.logger.Info(err)
		return "", customerr.NotAcceptable
	}

	return s.generateToken(objID)
}

func (s *Service) SignIn(ctx context.Context, username string, password string) (string, error) {
	user, err := s.userStorage.GetOne(ctx, username, s.generatePasswordHash(password))
	if err != nil {
		return "", customerr.NotFoundErr
	}

	return s.generateToken(user.ID)
}

func (s *Service) generateToken(id string) (string, error) {
	tkCl := tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tkCl)

	return token.SignedString([]byte(signKey))
}

func (s *Service) generatePasswordHash(password string) string {
	hash := sha256.New()

	_, err := hash.Write([]byte(password))
	if err != nil {
		s.logger.Info(err)
	}

	result := hash.Sum([]byte(salt))

	return fmt.Sprintf("%x", result)
}

func (s *Service) ParseToken(accessToken string) (string, error) {

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})

	if err != nil {
		return "", err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return "", fmt.Errorf("invalid signing method")
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", fmt.Errorf("token claims are not of internal type *tokenClaims")
	}

	return claims.UserID, nil
}

//func (s *Service) Update(ctx context.Context, user usermodel.UserInternal) error {
//	return s.userstorage.Update(ctx, user)
//}
//
//func (s *Service) Delete(ctx context.Context, id string) error {
//	return s.userstorage.Delete(ctx, id)
//}
