package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// this whole package is directly derived from my senior's repo (with some changes ofc)
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type JWTConfig struct {
	SecretKey          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
}

type Claims struct {
	ID        int64  `json:"user_id"`
	Email     string `json:"email"`
	IsAdmin   bool   `json:"is_admin"`
	TokenType string `json:"token_type"`

	jwt.RegisteredClaims
}

type tokenClaims struct {
	ID        int64
	Email     string
	IsAdmin   bool
	TokenType string
	Expiry    time.Duration
}

func DefaultJWTConfig(secretKey string) *JWTConfig {
	return &JWTConfig{
		SecretKey:          secretKey,
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		Issuer:             "go-ecommerce",
	}
}

func (c *JWTConfig) GenerateAccessToken(id int64, email string, isAdmin bool) (string, error) {
	return c.generateToken(tokenClaims{
		ID:        id,
		Email:     email,
		IsAdmin:   isAdmin,
		TokenType: "access",
		Expiry:    c.AccessTokenExpiry,
	})
}

func (c *JWTConfig) GenerateRefreshToken(id int64, email string, isAdmin bool) (string, error) {
	return c.generateToken(tokenClaims{
		ID:        id,
		Email:     email,
		IsAdmin:   isAdmin,
		TokenType: "refresh",
		Expiry:    c.RefreshTokenExpiry,
	})
}

func (c *JWTConfig) generateToken(tc tokenClaims) (string, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("generate token id: %w", err)
	}

	now := time.Now()

	claims := Claims{
		ID:        tc.ID,
		Email:     tc.Email,
		IsAdmin:   tc.IsAdmin,
		TokenType: tc.TokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Issuer:    c.Issuer,
			Subject:   tc.Email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tc.Expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(c.SecretKey))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

func (c *JWTConfig) ValidateToken(tokenString string) (*Claims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}

			return []byte(c.SecretKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}

		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
