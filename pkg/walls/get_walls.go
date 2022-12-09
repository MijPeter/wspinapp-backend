package walls

import (
	"example/wspinapp-backend/pkg/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetWalls(c *gin.Context) {
	var walls []common.Wall
	h.DB.Find(&walls)
	c.IndentedJSON(http.StatusOK, walls)
}
