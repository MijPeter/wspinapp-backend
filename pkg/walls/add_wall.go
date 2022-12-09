package walls

import (
	"example/wspinapp-backend/pkg/common"
	"example/wspinapp-backend/pkg/common/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func validateAddWall(err error) *errors.HttpError {
	if err != nil {
		return errors.BadRequest
	}

	return nil
}

func (h *handler) AddWall(c *gin.Context) {
	var newWall common.Wall

	err := c.BindJSON(&newWall) // this should be some other structure than common.Wall
	httpErr := validateAddWall(err)

	if httpErr != nil {
		c.IndentedJSON(httpErr.Status(), httpErr.Error())
		return
	}

	newWall.Id = strconv.Itoa(len(walls))
	walls = append(walls, newWall) // TODO use some kind of db or sth, maybe firestore for now and then as real db
	c.IndentedJSON(http.StatusCreated, newWall)
}
