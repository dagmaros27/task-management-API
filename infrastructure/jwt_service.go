package infrastructure

import (
	"fmt"
	"net/http"
	"task_managment_api/domain"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	GenerateUserToken(user domain.User) (string, domain.CustomError)
	ValidateToken(tokenString string) (jwt.MapClaims, domain.CustomError)
}

type jwtService struct{
	AccessTokenSecret string
}

func NewJWTService(secret string) JWTService {
	return &jwtService{
		AccessTokenSecret: secret,
	}
}


func (js *jwtService) GenerateUserToken(user domain.User) (string, domain.CustomError) {{

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &domain.Claims{
		UserId:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(js.AccessTokenSecret))
	if err != nil {
		return "", domain.CustomError{ErrCode: http.StatusInternalServerError,ErrMessage:  "Error while generating token"}
	}

	return tokenString, domain.CustomError{}
}}

func  (js *jwtService) ValidateToken(tokenString string) (jwt.MapClaims, domain.CustomError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(js.AccessTokenSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, domain.CustomError{ErrCode: http.StatusUnauthorized, ErrMessage: "Invalid token"}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.CustomError{ErrCode: http.StatusUnauthorized, ErrMessage: "Invalid token"}
	}

	return claims, domain.CustomError{}
}


