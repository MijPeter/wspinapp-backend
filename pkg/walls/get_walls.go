package walls

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetWalls(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, walls)

}
