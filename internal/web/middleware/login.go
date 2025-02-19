package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
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
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 0})
			return
		}
		id := sess.Get("userId")
		if id == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 0})
			return
		}
		userId, ok := id.(int64)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": 200, "code": 0})
			return
		}
		ctx.Set("user_id", userId)

	}
}
