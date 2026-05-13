package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// PasswordHash 对密码进行加密
func PasswordHash(password string) (string, error) {
	// GenerateFromPassword 方法会自动加盐，bcrypt.DefaultCost 是推荐的加强强度（10）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// PasswordVerify 验证密码是否正确
// 参数：
//
// hashedPassword - 数据库中存储的加密密码
// password - 用户输入的密码
//
// 返回 true 密码正确，false 密码错误
func PasswordVerify(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// JWTClaims JWT 载荷信息
type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"` // user / creator / admin
	jwt.RegisteredClaims
}

func GenerateToken(userID int64, username, role, secret string, expire int64) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间 = 当前时间 + expire 秒
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expire) * time.Second)),
			// 签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// 生效时间（立即生效）
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 使用 HS256 算法生成 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并生成 Token 字符串
	return token.SignedString([]byte(secret))
}

// ParseToken 解析 JWT Token
// 参数：
//
// tokenString - Token 字符串
// secret - 密钥（必须与生成时一致）
//
// 返回: Claims 和可能的错误
func ParseToken(tokenString, secret string) (*JWTClaims, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 返回密钥用于验证签名
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	// 验证 JWT
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
