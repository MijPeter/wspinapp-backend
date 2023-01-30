package controller

import "github.com/gin-gonic/gin"

func (h *routesHandler) Pong(c *gin.Context) {
	c.String(200, "pong")
}
