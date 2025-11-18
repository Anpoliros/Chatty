package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT claims结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTUtil JWT工具
type JWTUtil struct {
	secret     []byte
	expireTime time.Duration
}

// NewJWTUtil 创建JWT工具
func NewJWTUtil(secret string, expireHours int) *JWTUtil {
	return &JWTUtil{
		secret:     []byte(secret),
		expireTime: time.Duration(expireHours) * time.Hour,
	}
}

// GenerateToken 生成JWT token
func (j *JWTUtil) GenerateToken(userID uint, username string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ParseToken 解析JWT token
func (j *JWTUtil) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新token
func (j *JWTUtil) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 如果token即将过期(剩余时间小于1小时),则刷新
	if time.Until(claims.ExpiresAt.Time) < time.Hour {
		return j.GenerateToken(claims.UserID, claims.Username)
	}

	return tokenString, nil
}
