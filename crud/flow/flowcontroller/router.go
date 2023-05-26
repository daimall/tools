package flowcontroller

import (
	"strings"

	"github.com/gin-gonic/gin"
)

var fc *FlowController

func init() {
	if fc == nil {
		fc = &FlowController{}
		fc.Mapping("Post", fc.Post)
		fc.Mapping("GetOne", fc.GetOne)
		fc.Mapping("GetAll", fc.GetAll)
		fc.Mapping("Put", fc.Put)
		fc.Mapping("Action", fc.Action)
		fc.Mapping("Delete", fc.Delete)
		fc.Mapping("DeleteList", fc.DeleteList)
	}
}

/*
路由注册，样例如下
Router("/v1/crud/:service", "get:GetAll;post:Post", gin.HandlerFunc...)
Router("/v1/crud/:service/:id:int", "get:GetOne;put:Put;delete:Delete", gin.HandlerFunc...)
Router("/v1/crud/:service/deletelist", "delete:DeleteList", gin.HandlerFunc...)
Router("/v1/crud/:service/:id/:action", "post:Action;get:Action", gin.HandlerFunc...)
*/
func Router(r *gin.Engine, rootPath string, mappingMethods string, handlers ...gin.HandlerFunc) {
	// 启动session
	for _, kv := range strings.Split(mappingMethods, ";") {
		methods := strings.Split(kv, ":")
		if len(methods) != 2 {
			panic("mappingMethods[ " + kv + " ] invalid")
		}
		mk, mv := methods[0], methods[1]
		hs := []gin.HandlerFunc{fc.GetMapMethod(mv)}
		hs = append(hs, handlers...)
		switch mk {
		case "get":
			r.GET(rootPath, hs...)
		case "post":
			r.POST(rootPath, hs...)
		case "put":
			r.PUT(rootPath, hs...)
		case "patch":
			r.PATCH(rootPath, hs...)
		case "delete":
			r.DELETE(rootPath, hs...)
		case "head":
			r.HEAD(rootPath, hs...)
		case "options":
			r.OPTIONS(rootPath, hs...)
		default:
			panic("http methed[ " + mk + " ] not support")
		}
	}
}
