package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{paths: make([]string, 0)}
}

func (l *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	//gob.Register(time.Time{}) // 如果使用 Redis 缓存 time.time,需要先注册格式,要不然读取会出现异常
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		ignore := false
		for _, v := range l.paths {
			if v == path {
				ignore = true
			}
		}
		if ignore {
			return
		}
		sess := sessions.Default(ctx)
		if sess == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 50010, "msg": "用户未登录"})
			return
		}
		id := sess.Get("userId")
		if id == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 50010, "msg": "用户未登录"})
			return
		}
		userId, ok := id.(int64)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 50010, "msg": "用户未登录"})
			return
		}
		ctx.Set("user_id", userId)

		// 刷新 session
		now := time.Now().UnixMilli()
		updateTime := sess.Get("update_time")
		if updateTime == nil {
			sess.Set("userId", userId)
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		updateTimeVal, _ := updateTime.(int64)
		if now-updateTimeVal > 60*1000 {
			sess.Set("update_time", now)
			sess.Set("userId", userId)
			sess.Save()
			return
		}

	}
}
