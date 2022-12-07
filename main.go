package main

import (
	"example/wspinapp-backend/db"
	"example/wspinapp-backend/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	{Id: "0", Holds: []hold{}, Image: "image_1.png"},
}

var routes []route
var basicAuth = gin.BasicAuth(gin.Accounts{
	"wspinapp": "wspinapp",
})

func main() {
	db.ConnectDb()

	router := gin.Default()
	router.Use(basicAuth)
	router.GET("/walls", getWalls)
	router.POST("/walls", addWall)
	router.GET("/walls/:wallId", getWall)
	router.GET("/walls/:wallId/routes", getRoutes)

	router.Run()
}

// TODO validation should be probably done some other way
func validateAddWall(err error) *errors.HttpError {
	if err != nil {
		return errors.BadRequest
	}

	return nil
}

func getWalls(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, walls)
}

func addWall(c *gin.Context) {
	var newWall wall

	err := c.BindJSON(&newWall)
	httpErr := validateAddWall(err)

	if httpErr != nil {
		c.IndentedJSON(httpErr.Status(), httpErr.Error())
		return
	}

	newWall.Id = strconv.Itoa(len(walls))
	walls = append(walls, newWall) // TODO use some kind of db or sth, maybe firestore for now and then as real db
	c.IndentedJSON(http.StatusCreated, newWall)
}

func getWall(c *gin.Context) {
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

func getRoutes(c *gin.Context) {
	wallId := c.Param("wallId")

	wallRoutes := make([]route, 0)
	for _, r := range routes {
		if r.Wall.Id == wallId {
			wallRoutes = append(wallRoutes, r)
		}

	}
	c.IndentedJSON(http.StatusOK, wallRoutes)
}
