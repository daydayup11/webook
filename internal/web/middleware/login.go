package middleware

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/m/internal/web"
)

type LoginMiddlewareBuilder struct {
	path []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}
func (lb *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	lb.path = append(lb.path, path)
	return lb
}

func (lb *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(context *gin.Context) {
		for _, path := range lb.path {
			if path == context.Request.URL.Path {
				return
			}
		}
		session := sessions.Default(context)
		userId := session.Get(web.UserIdKey)
		if userId == nil {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

//var IgnorePaths []string

//func (lb *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
//	return func(context *gin.Context) {
//		for _, path := range lb.path {
//			if path == context.Request.URL.Path {
//				return
//			}
//		}
//
//		session := sessions.Default(context)
//
//	}
//}
