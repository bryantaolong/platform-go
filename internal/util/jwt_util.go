package util

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
	"time"
)

type JWTUtil struct {
	secretKey  []byte
	expiration time.Duration
}

func NewJWTUtil(secret string) *JWTUtil {
	return &JWTUtil{
		secretKey:  []byte(secret),
		expiration: 24 * time.Hour, // 24小时过期
	}
}

type Claims struct {
	UserID uint64 `json:"user_id"`
	Roles  string `json:"roles"`
	jwt.RegisteredClaims
}

func (j *JWTUtil) GenerateToken(userID uint64, roles string) (string, error) {
	claims := Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(userID), 10),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTUtil) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (j *JWTUtil) GetRolesFromToken(tokenString string) ([]string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Roles == "" {
		return []string{}, nil
	}

	roles := strings.Split(claims.Roles, ",")
	var processedRoles []string

	for _, role := range roles {
		role = strings.TrimSpace(role)
		if !strings.HasPrefix(role, "ROLE_") {
			role = "ROLE_" + role
		}
		processedRoles = append(processedRoles, role)
	}

	return processedRoles, nil
}
