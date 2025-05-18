package pkg

import (
	"strconv"
	"time"

	"github.com/YangZhaoWeblog/UserService/internal/conf"
	"github.com/golang-jwt/jwt/v4" // 使用v4版本
)

type CustomClaims struct {
	UserID               string `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // 使用RegisteredClaims替代StandardClaims
}

type JwtClient struct {
	signingKey  string
	expiresTime int32
}

// NewClient 创建一个新的JWT客户端
func NewClient(cfg *conf.Data) *JwtClient {
	return &JwtClient{
		signingKey:  cfg.Jwt.SigningKey,
		expiresTime: cfg.Jwt.ExpiresTime,
	}
}

// GenerateToken 生成访问令牌和刷新令牌
func (c *JwtClient) GenerateToken(userID int64, username string) (accessToken string, refreshToken string, err error) {
	// 生成访问令牌
	now := time.Now()
	expireTime := now.Add(time.Duration(c.expiresTime) * time.Second)

	accessClaims := CustomClaims{
		UserID:   strconv.FormatInt(userID, 10),
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = token.SignedString([]byte(c.signingKey))
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshExpireTime := now.Add(time.Duration(c.expiresTime*7) * time.Second)
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(refreshExpireTime),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = token.SignedString([]byte(c.signingKey))

	return accessToken, refreshToken, err
}

// ParseToken 解析令牌
func (c *JwtClient) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.signingKey), nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
