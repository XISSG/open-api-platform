package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const SecretJWTKey = "dgrijalva"

var JWTExpireTime = time.Now().Add(time.Hour * 10)

func GenerateJWT(user string, role string, expire time.Time, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"role": role,
		"exp":  expire,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

type JWTResult struct {
	User      string
	Role      string
	Exp       time.Time
	Signature string
	Err       error
}

func ParseJWT(tokenString string, secret []byte) *JWTResult {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证 token 使用的签名算法是否符合预期
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return &JWTResult{Err: err}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return &JWTResult{Err: err}
	}

	return &JWTResult{
		User:      claims["user"].(string),
		Role:      claims["role"].(string),
		Exp:       time.Unix(int64(claims["exp"].(float64)), 0),
		Signature: token.Signature,
		Err:       nil,
	}
}
