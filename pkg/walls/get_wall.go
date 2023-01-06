package walls

import (
	"example/wspinapp-backend/pkg/common"
	"example/wspinapp-backend/pkg/common/errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *routesHandler) GetWall(c *gin.Context) {
	wallId64, err := strconv.ParseUint(c.Param("wallId"), 10, 32)
	wallId := uint(wallId64)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	var wall common.Wall

	err = h.database.First(&wall, wallId).Error

	if err != nil {
		httpErr := errors.NotFound
		c.IndentedJSON(httpErr.Status(), httpErr.Error())
		return
	}

	log.Printf("Found wall with id: %d\n", wallId)
	c.IndentedJSON(http.StatusOK, wall)
}
