package walls

import (
	"example/wspinapp-backend/pkg/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *handler) GetRoutes(c *gin.Context) {
	wallId64, err := strconv.ParseUint(c.Param("wallId"), 10, 32)
	wallId := uint(wallId64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}

	var wallRoutes []common.Route
	h.DB.Where(common.Route{WallID: wallId}).Find(&wallRoutes)

	c.IndentedJSON(http.StatusOK, wallRoutes)
}
