package walls

import (
	"example/wspinapp-backend/pkg/common/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) GetWall(c *gin.Context) {
	wallId := c.Param("wallId")
	for _, w := range walls {
		if w.Id == wallId {
			c.IndentedJSON(http.StatusOK, w)
			return
		}
	}

	httpErr := errors.NotFound
	c.IndentedJSON(httpErr.Status(), httpErr.Error())
}
