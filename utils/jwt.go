package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWT struct {
	user   string
	exp    time.Time
	secret string
}

func (j *JWT) GenerateJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": j.user,
		"exp":  j.exp,
	})

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证 token 使用的签名算法是否符合预期
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Username:", claims["user"])
		fmt.Println("Expires at:", claims["exp"])
	} else {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
