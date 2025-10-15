package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"jike_gin/consts"
	"jike_gin/internal/web"
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
		// webook前端 里面拼接了 Bearer
		authorization := ctx.GetHeader("authorization")
		args := strings.Split(authorization, " ")
		if len(args) != 2 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 50010, "msg": "长度问题:用户未登录"})
			zap.L().Info("Build", zap.String("error", "长度不够"))
			return
		}
		tokenStr := args[1]

		userClaims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, userClaims, func(token *jwt.Token) (i interface{}, e error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), e
		})
		if err != nil {
			fmt.Println("Build ParseWithClaims error", err.Error())
			zap.L().Info("Build ParseWithClaims error", zap.Error(err), zap.String("token", tokenStr))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if token == nil || !token.Valid {
			fmt.Println("Build Valid error")
			zap.L().Info("Build Valid error", zap.Error(err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//使用 registered Claims的 expired_at ,解析不要判断过期时间,只判断是否有效 token.Valid
		// 确定 auth 是否有效
		//claims, _ := token.Claims.(jwt.MapClaims)
		//if lts, ok := claims["expire_at"].(string); !ok || time.Now().UTC().Format(time.DateTime) > lts {
		//	err = fmt.Errorf("token expired")
		//	fmt.Println("build_jwt===>", err, claims)
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}

		if userClaims.UserAgent != ctx.Request.UserAgent() {
			// 风险检查,设备信息是否一致
			//ctx.AbortWithStatus(http.StatusUnauthorized)

			zap.L().Info("Build agent error", zap.String("reqAgent", ctx.Request.UserAgent()), zap.String("userAgent", userClaims.UserAgent))
			//return
		}

		userClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
		tokenStr, _ = token.SignedString(consts.GetAuthSecret())
		// 续约
		ctx.Header("x-jwt-token", tokenStr)
		fmt.Println("build_jwt uid===>", userClaims.UId)
		ctx.Set("user_id", userClaims.UId)
	}
}
