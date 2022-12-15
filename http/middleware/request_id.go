package middleware

import (
	"github.com/bigtree8/flow-core/utils"
	"github.com/gin-gonic/gin"
)

func GenRequestId(c *gin.Context) {
	reqId := utils.GenUUID()
	c.Request.Header.Add("X-Request-Id", reqId)
	c.Header("X-Request-Id", reqId)
	c.Next()
}
