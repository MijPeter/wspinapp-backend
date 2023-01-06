package walls

import (
	"example/wspinapp-backend/pkg/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *routesHandler) GetWalls(c *gin.Context) {
	var walls []common.Wall
	h.database.Find(&walls)
	c.IndentedJSON(http.StatusOK, walls)
}
