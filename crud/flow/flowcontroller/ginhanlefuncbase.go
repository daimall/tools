package flowcontroller

import (
	"github.com/gin-gonic/gin"
)

type FuncAdapter interface {
}

func HandlerFuncBase(c *gin.Context, adapter FuncAdapter) {

}
