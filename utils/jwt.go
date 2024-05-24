package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const SecretJWTKey = "dgrijalva"

var JWTExpireTime = time.Now().Add(time.Hour * 10)

func GenerateJWT(user string, role string, expire int64, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        user,
		ExpiresAt: expire,
		Audience:  role,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseJWT(tokenString string, secret []byte) (*jwt.StandardClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证 token 使用的签名算法是否符合预期
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenClaims == nil {
		return nil, err
	}

	claims, ok := tokenClaims.Claims.(*jwt.StandardClaims)
	if ok && tokenClaims.Valid {
		return claims, nil
	}
	return nil, err
}
