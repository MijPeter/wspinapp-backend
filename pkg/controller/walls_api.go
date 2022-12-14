package controller

import (
	"example/wspinapp-backend/pkg/common"
	"example/wspinapp-backend/pkg/common/errors"
	"example/wspinapp-backend/pkg/services/walls_service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *routesHandler) AddWall(c *gin.Context) {
	var newWall common.Wall

	err := c.BindJSON(&newWall)

	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}
	walls_service.AddWall(h.database, &newWall)
	c.IndentedJSON(http.StatusCreated, newWall)
}

func (h *routesHandler) GetWall(c *gin.Context) {
	wallId64, err := strconv.ParseUint(c.Param("wallId"), 10, 32)
	wallId := uint(wallId64)

	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	wall, err := walls_service.GetWall(h.database, wallId)

	if err != nil {
		returnErrorResponse(c, errors.NotFound)
		return
	}

	log.Printf("Found wall with id: %d\n", wallId)
	c.IndentedJSON(http.StatusOK, wall)
}

func (h *routesHandler) GetWalls(c *gin.Context) {
	walls := walls_service.GetWalls(h.database)
	c.IndentedJSON(http.StatusOK, walls)
}

func (h *routesHandler) GetRoutes(c *gin.Context) {
	wallId64, err := parseUint(c.Param("wallId"))
	wallId := uint(wallId64)
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	wallRoutes := walls_service.GetRoutes(h.database, wallId)

	c.IndentedJSON(http.StatusOK, wallRoutes)
}

func (h *routesHandler) UploadImage(c *gin.Context) {
	uploadedFile, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Failed to parse given image, %s\n", err.Error())
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	//validate
	err = h.validate.Struct(uploadedFile)
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	wallId64, err := parseUint(c.Param("wallId"))
	wallId := uint(wallId64)

	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	uploadUrl, err := walls_service.UploadFileAndSaveUrlToDb(
		h.database,
		h.imageRepository,
		wallId,
		common.File{File: uploadedFile})

	if err != nil {
		log.Printf("Failed to upload image, %s\n", err.Error())
		returnErrorResponse(c, errors.InternalError)
		return
	}

	c.IndentedJSON(http.StatusCreated, uploadUrl)
}

func returnErrorResponse(c *gin.Context, httpError *errors.HttpError) {
	c.IndentedJSON(
		httpError.Status(),
		httpError.Error())
}

func parseUint(id string) (uint64, error) {
	return strconv.ParseUint(id, 10, 32)
}

//https://dev.to/hackmamba/robust-media-upload-with-golang-and-cloudinary-gin-gonic-version-54ii
