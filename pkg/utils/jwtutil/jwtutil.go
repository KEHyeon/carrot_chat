// jwtutil/jwtutil.go
package jwtutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config 구조체는 JWT 패키지의 설정값을 담습니다.
type Config struct {
	SecretKey      string        // 서명 키
	ExpireDuration time.Duration // 토큰 만료 시간
}

// Claims 구조체는 기본 JWT 클레임을 담습 니다.
type Claims struct {
	UserID uint64 `json:"user_id"` // 유저 ID
	jwt.RegisteredClaims
}

// JWTUtil는 JWT 생성 및 검증을 위한 유틸리티입니다.
type JWTUtil struct {
	config Config
}

// NewJWTUtil은 JWTUtil을 초기화합니다.
func NewJWTUtil(secretKey string, expireDuration time.Duration) *JWTUtil {
	return &JWTUtil{
		config: Config{
			SecretKey:      secretKey,
			ExpireDuration: expireDuration,
		},
	}
}

// GenerateToken은 유저 ID를 기반으로 JWT 토큰을 생성합니다.
func (j *JWTUtil) GenerateToken(userID uint64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.config.ExpireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.config.SecretKey))
}

// ValidateToken은 JWT 토큰을 검증하고 클레임을 반환합니다.
func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
