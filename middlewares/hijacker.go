package middlewares

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

type (
	// CookieHacker 将cookie转换为http header
	CookieHacker struct {
		CookieName string `json:"cookieName"` // 原始的cookie名称
		HeaderName string `json:"headerName"` // 目标的header名称

		name string
	}

	bodyWriter struct {
		gin.ResponseWriter
		body       *bytes.Buffer
		headerName string
	}
)

var (
	// DefaultCookieHacker 默认的cookie转换header功能
	DefaultCookieHacker = NewCookieHacker("cookieHacker")
)

// Name 配置文件名称
func (c *CookieHacker) Name() string {
	return c.name
}

// Hijacker support session id by header
func (c *CookieHacker) Hijacker() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 如果cookie存在，那么直接返回了，不进行hack
		cookie, _ := ctx.Cookie(c.CookieName)
		if cookie != "" {
			ctx.Next()
			return
		}

		// 请求header中有token存在，那么设置到请求的cookie中方便将session解析出来
		token := ctx.GetHeader(c.HeaderName)
		if token != "" {
			ctx.Request.Header.Add("Cookie", token)
		}

		// hack ResponseWriter
		bodyWriter := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer, headerName: c.HeaderName}
		ctx.Writer = bodyWriter

		ctx.Next()
	}
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)

	setCookieValue := w.Header().Get("set-cookie")
	if setCookieValue != "" {
		w.Header().Add(w.headerName, setCookieValue)
	}

	return w.ResponseWriter.Write(b)
}

// NewCookieHacker 返回一个CookieHacker
func NewCookieHacker(name string) *CookieHacker {
	return &CookieHacker{
		name: name,
	}
}
