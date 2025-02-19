package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go_project/gin/consts"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{paths: make([]string, 0)}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePath(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Time{}) // 如果使用 Redis 缓存 time.time,需要先注册格式,要不然读取会出现异常
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		for _, v := range l.paths {
			if v == path {
				return
			}
		}
		// webook 里面拼接了 banner
		authorization := ctx.GetHeader("authorization")
		args := strings.Split(authorization, " ")
		if len(args) != 2 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 50010, "msg": "用户未登录"})
			return
		}
		tokenStr := args[1]
		fmt.Println("token", tokenStr)

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, e error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(consts.GetAuthSecret()), e
		})
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 50010, "msg": "用户未登录"})
			return
		}
		fmt.Println("token.Valid", token.Valid)
		if token == nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 50010, "msg": "用户未登录"})
			return
		}

		// 确定 auth 是否有效
		claims, _ := token.Claims.(jwt.MapClaims)
		if lts, ok := claims["expire_at"].(string); !ok || time.Now().UTC().Format(time.DateTime) > lts {
			err = fmt.Errorf("token expired")
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 50010, "msg": "用户未登录"})
			return
		}
		fmt.Println("claims", claims)
		fmt.Println("id", claims["id"])

		//解析用户的 ID
		id, ok := claims["id"]
		if !ok {
			err = fmt.Errorf("not id")
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 50010, "msg": "用户未登录"})
			return
		}

		userId, ok := id.(float64)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 50010, "msg": "用户未登录"})
			return
		}
		ctx.Set("user_id", int64(userId))

	}
}
