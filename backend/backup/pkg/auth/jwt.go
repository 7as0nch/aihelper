package auth

import (
	"errors"
	"github.com/aichat/backend/models"
	"github.com/aichat/backend/pkg/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var jwtSecret []byte

// 初始化JWT配置
func InitJWT() {
	jwtSecret = []byte(viper.GetString("jwt.secret"))
}

// CustomClaims 自定义Claims
type CustomClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT
func GenerateToken(userID uint64, username string, role string) (string, error) {
	expireTime := time.Now().Add(time.Duration(viper.GetInt64("jwt.expire_time")) * time.Hour)

	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "aichat",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetUserIDFromToken 从token中获取用户ID
func GetUserIDFromToken(c *gin.Context) (uint64, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, errors.New("authorization header is empty")
	}

	tokenString := authHeader[len("Bearer "):]
	claims, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}

// GetUserInfoFromToken 从token中获取用户信息
func GetUserInfoFromToken(c *gin.Context) (models.User, error) {
	var user models.User
	userID, err := GetUserIDFromToken(c)
	if err != nil {
		return user, err
	}

	db := db.GetDB()
	db.First(&user, userID)
	if user.ID == 0 {
		return user, errors.New("user not found")
	}

	return user, nil
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 提取token
		tokenString := authHeader[len("Bearer "):]
		claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 检查token是否过期
		if time.Now().Unix() > claims.ExpiresAt {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// 将用户信息保存到上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware 管理员中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户角色
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin permission required"})
			c.Abort()
			return
		}

		c.Next()
	}
}