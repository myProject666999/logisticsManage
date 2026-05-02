package utils

import (
	"time"

	"logisticsManage/config"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"` // user, delivery, admin
	jwt.StandardClaims
}

// 生成JWT Token
func GenerateToken(userID uint, username string, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireTime) * time.Second).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "logisticsManage",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.NewValidationError("无效的Token", jwt.ValidationErrorClaimsInvalid)
}
