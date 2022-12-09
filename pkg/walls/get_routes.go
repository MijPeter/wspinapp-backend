package walls

import (
	"example/wspinapp-backend/pkg/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetRoutes(c *gin.Context) {
	wallId := c.Param("wallId")

	wallRoutes := make([]common.Route, 0)
	for _, r := range routes {
		if r.WallId == wallId {
			wallRoutes = append(wallRoutes, r)
		}

	}
	c.IndentedJSON(http.StatusOK, wallRoutes)
}
