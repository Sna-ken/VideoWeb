package jwt

import (
	"errors"
	"time"

	"github.com/Sna-ken/videoweb/config"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateAccessToken(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Second * time.Duration(config.JWT.AccessTokenExpiry)).Unix() //token有效时间
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Issuer:    "Snaken-Video-Web",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //用签名算法创建一个新token

	signerToken, err := token.SignedString([]byte(config.JWT.AccessTokenSecret)) //使用secretKey签名
	if err != nil {
		return "", err
	}
	return signerToken, err
}

func GenerateRefreshToken(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Second * time.Duration(config.JWT.RefreshTokenExpiry)).Unix() //token有效时间
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Issuer:    "Snaken-Video-Web",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //用签名算法创建一个新token

	signerToken, err := token.SignedString([]byte(config.JWT.RefreshTokenSecret)) //使用secretKey签名
	if err != nil {
		return "", err
	}
	return signerToken, err
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.JWT.AccessTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims) //获取解析后的Claims
	if !ok || !token.Valid {
		return nil, errors.New("claims invalid")
	}

	return claims, nil
}

func ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.JWT.RefreshTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims) //获取解析后的Claims
	if !ok || !token.Valid {
		return nil, errors.New("claims invalid")
	}

	return claims, nil
}
