package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	GenerateToken(userId int) (string, error)
}

type jwtService struct {
}

func NewService() *jwtService {
	return &jwtService{}
}

var SECRET_KEY = []byte("asdasdasdasdasd")

func (s *jwtService) GenerateToken(userId int) (string, error) {

	claim := jwt.MapClaims{}
	claim["user_id"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedString, err
	}

	return signedString, nil

}
