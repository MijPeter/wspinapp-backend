package main

import (
	"example/wspinapp-backend/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type hold struct {
	Id       string   `json:"id"`
	Position position `json:"position"`
}

type wall struct {
	Id    string `json:"id"`
	Holds []hold `json:"holds"`
	Image string `json:"image"`
}

type route struct {
	Id    string `json:"id"`
	Holds []hold `json:"holds"`
	Wall  wall   `json:"wall"`
}

var walls = []wall{
	{Id: "1", Holds: []hold{}, Image: "image_1.png"},
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil) // this should be revisited
	router.GET("/walls", getWalls)
	router.POST("/walls", addWall)
	router.GET("/walls/:id", getWall)

	router.Run("localhost:8080")
}

// inner method that should be in an inner layer (touching db)
func validateAddWall(err error, newWall wall) *errors.HttpError {
	if err != nil {
		return errors.BadRequest
	}

	for _, existingWall := range walls {
		if existingWall.Id == newWall.Id {
			return errors.Conflict
		}
	}
	return nil
}

func getWalls(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, walls)
}

func addWall(c *gin.Context) {
	var newWall wall

	err := c.BindJSON(&newWall)
	httpErr := validateAddWall(err, newWall)

	if httpErr != nil {
		c.IndentedJSON(httpErr.Status(), httpErr.Error())
		return
	}

	walls = append(walls, newWall) // TODO use some kind of db or sth, maybe firestore for now and then as real db
	c.IndentedJSON(http.StatusCreated, newWall)
}

func getWall(c *gin.Context) {
	id := c.Param("id")

	for _, w := range walls {
		if w.Id == id {
			c.IndentedJSON(http.StatusOK, w)
			return
		}
	}

	err := errors.NotFound
	c.IndentedJSON(err.Status(), err.Error())
}
