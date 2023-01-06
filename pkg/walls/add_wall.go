package walls

import (
	"example/wspinapp-backend/pkg/common"
	"example/wspinapp-backend/pkg/common/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func validateAddWall(err error) *errors.HttpError {
	if err != nil {
		return errors.BadRequest
	}

	return nil
}

func (h *routesHandler) AddWall(c *gin.Context) {
	var newWall common.Wall

	err := c.BindJSON(&newWall) // this should be some other structure than common.Wall
	httpErr := validateAddWall(err)

	if httpErr != nil {
		c.IndentedJSON(httpErr.Status(), httpErr.Error())
		return
	}
	h.database.Create(&newWall)
	c.IndentedJSON(http.StatusCreated, newWall)
}
