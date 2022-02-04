package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/golang-jwt/jwt"
	"nprn/internal/customerr"
	"nprn/internal/entity/sale/salemodel"
	"nprn/internal/entity/user/usermodel"
	"nprn/pkg/logging"
	"time"
)

const (
	salt      = "4hsd83jd7fsd2"
	tokenTime = 12 * time.Hour
	signKey   = "dkr3!#mc349x#s3&74f12d"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type SaleStorage interface {
	Create(ctx context.Context, sale salemodel.Sale) (string, error)
	GetOne(ctx context.Context, id string) (salemodel.Sale, error)
	GetAll(ctx context.Context) ([]salemodel.Sale, error)
	Update(ctx context.Context, sale salemodel.Sale) error
	Delete(ctx context.Context, id string) error
}

type UserStorage interface {
	Create(ctx context.Context, user usermodel.UserInternal) (string, error)
	GetOne(ctx context.Context, username string, password string) (usermodel.UserTransfer, error)
	//Update(ctx context.Context, user usermodel.UserInternal) error
	//Delete(ctx context.Context, id string) error
}

type Service struct {
	UserStorage UserStorage
	SaleStorage SaleStorage
	Logger      *logging.Logger
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

func NewService(userStorage UserStorage, saleStorage SaleStorage, logger *logging.Logger) *Service {
	return &Service{
		UserStorage: userStorage,
		SaleStorage: saleStorage,
		Logger:      logger,
	}
}

func (s *Service) SignUp(ctx context.Context, user usermodel.UserInternal) (string, error) {
	passHash, err := GeneratePasswordHash(user.PasswordHash)
	if err != nil {
		return "", err
	}

	user.PasswordHash = passHash

	objID, err := s.UserStorage.Create(ctx, user)
	if err != nil {
		s.Logger.Info(err)
		return "", customerr.NotAcceptable
	}

	return GenerateToken(objID)
}

func (s *Service) SignIn(ctx context.Context, username string, password string) (string, error) {
	passHash, err := GeneratePasswordHash(password)
	if err != nil {
		return "", err
	}

	user, err := s.UserStorage.GetOne(ctx, username, passHash)
	if err != nil {
		return "", customerr.NotFoundErr
	}

	return GenerateToken(user.ID)
}

func GenerateToken(id string) (string, error) {
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

func GeneratePasswordHash(password string) (string, error) {
	hash := sha256.New()

	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", err
	}

	result := hash.Sum([]byte(salt))

	return fmt.Sprintf("%x", result), nil
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

//func (s *Service) UpdateUser(ctx context.Context, user usermodel.UserInternal) error {
//	return s.userstorage.Update(ctx, user)
//}
//
//func (s *Service) DeleteUser(ctx context.Context, id string) error {
//	return s.userstorage.Delete(ctx, id)
//}

func (s *Service) CreateSale(ctx context.Context, sale salemodel.Sale) (string, error) {
	return s.SaleStorage.Create(ctx, sale)
}

func (s *Service) GetSale(ctx context.Context, id string) (salemodel.Sale, error) {
	return s.SaleStorage.GetOne(ctx, id)
}

func (s *Service) GetAllSales(ctx context.Context) ([]salemodel.Sale, error) {
	return s.SaleStorage.GetAll(ctx)
}

func (s *Service) UpdateSale(ctx context.Context, sale salemodel.Sale) error {
	return s.SaleStorage.Update(ctx, sale)
}

func (s *Service) DeleteSale(ctx context.Context, id string) error {
	return s.SaleStorage.Delete(ctx, id)
}
