package flowcontroller

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func UseMiddleWare(r *gin.Engine) {
	store := sessions.NewCookieStore([]byte("secret"))
	// 配置Session中间件
	r.Use(sessions.Sessions("session", store))
	r.Use(PrepareMiddleWare())
	r.Use(ResponseJSONMiddleWare())
}
