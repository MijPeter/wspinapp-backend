package controller

import (
	"example/wspinapp-backend/pkg/common/errors"
	"example/wspinapp-backend/pkg/common/schema"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *routesHandler) AddWall(c *gin.Context) {
	var newWall schema.Wall

	err := c.BindJSON(&newWall)

	if err != nil {
		returnErrorResponseDebug(c, errors.BadRequest, err)
		return
	}
	h.service.WallsService.AddWall(&newWall)
	c.IndentedJSON(http.StatusCreated, newWall)
}

func (h *routesHandler) GetWall(c *gin.Context) {
	wallId64, err := parseUint(c.Param("wallId"))
	wallId := uint(wallId64)

	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	wall, err := h.service.WallsService.GetWall(wallId)

	if err != nil {
		returnErrorResponse(c, errors.NotFound)
		return
	}

	log.Printf("Found wall with id: %d\n", wallId)
	c.IndentedJSON(http.StatusOK, wall)
}

func (h *routesHandler) UpdateWall(c *gin.Context) {
	wallId64, err := parseUint(c.Param("wallId"))
	wallId := uint(wallId64)

	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	var wall schema.Wall

	err = c.BindJSON(&wall)

	if err != nil {
		returnErrorResponseDebug(c, errors.BadRequest, err)
		return
	}

	stateWall, err := h.service.WallsService.UpdateWall(wallId, &wall)

	if err != nil {
		returnErrorResponseDebug(c, errors.BadRequest, err)
		return
	}

	c.IndentedJSON(http.StatusOK, stateWall)
}

func (h *routesHandler) DeleteWall(c *gin.Context) {
	wallId64, err := parseUint(c.Param("wallId"))
	wallId := uint(wallId64)

	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	err = h.service.WallsService.DeleteWall(wallId)
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}
	c.IndentedJSON(http.StatusNoContent, nil)
}

func (h *routesHandler) GetWalls(c *gin.Context) {
	walls := h.service.WallsService.GetWalls()
	c.IndentedJSON(http.StatusOK, walls)
}

func (h *routesHandler) GetRoutes(c *gin.Context) {
	wallId64, err := parseUint(c.Param("wallId"))
	wallId := uint(wallId64)
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	wallRoutes := h.service.WallsService.GetRoutes(wallId)

	c.IndentedJSON(http.StatusOK, wallRoutes)
}

func (h *routesHandler) AddRoute(c *gin.Context) {
	wallId64, err := parseUint(c.Param("wallId"))
	wallId := uint(wallId64)
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	var newRoute schema.Route
	err = c.BindJSON(&newRoute)

	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	err = h.service.WallsService.AddRoute(&newRoute, wallId)
	if err != nil {
		log.Printf(err.Error())
		returnErrorResponse(c, errors.BadRequest)
		return
	}
	c.IndentedJSON(http.StatusCreated, newRoute)
}

func (h *routesHandler) GetRoute(c *gin.Context) {
	wallId64, err := parseUint(c.Param("wallId"))
	wallId := uint(wallId64)
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	routeId64, err := parseUint(c.Param("routeId"))
	routeId := uint(routeId64)
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	route, err := h.service.WallsService.GetRoute(wallId, routeId)
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}
	c.IndentedJSON(http.StatusOK, route)
}

func (h *routesHandler) UpdateRoute(c *gin.Context) {
	wallId64, err := parseUint(c.Param("wallId"))
	wallId := uint(wallId64)
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	routeId64, err := parseUint(c.Param("routeId"))
	routeId := uint(routeId64)
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	var newRoute schema.Route
	err = c.BindJSON(&newRoute)

	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	updatedRoute, err := h.service.WallsService.UpdateRoute(&newRoute, wallId, routeId)
	if err != nil {
		log.Printf(err.Error())
		returnErrorResponse(c, errors.BadRequest)
		return
	}
	c.IndentedJSON(http.StatusOK, updatedRoute)
}

func (h *routesHandler) DeleteRoute(c *gin.Context) {
	wallId64, err := parseUint(c.Param("wallId"))
	wallId := uint(wallId64)
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	routeId64, err := parseUint(c.Param("routeId"))
	routeId := uint(routeId64)
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		return
	}

	err = h.service.WallsService.DeleteRoute(wallId, routeId)
	if err != nil {
		log.Printf(err.Error())
		returnErrorResponse(c, errors.BadRequest)
		return
	}
	c.IndentedJSON(http.StatusNoContent, "")
}

func (h *routesHandler) UploadImageFull(c *gin.Context) {
	h.UploadImage(c, false)
}

func (h *routesHandler) UploadImagePreview(c *gin.Context) {
	h.UploadImage(c, true)
}

func (h *routesHandler) UploadImage(c *gin.Context, isPreview bool) {
	uploadedFile, _, err := c.Request.FormFile("file")
	if err != nil {
		returnErrorResponse(c, errors.BadRequest)
		returnErrorResponseDebug(c, errors.BadRequest, err)
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

	uploadUrl, err := h.service.WallsService.UploadFileAndSaveUrlToDb(
		wallId,
		schema.File{File: uploadedFile},
		isPreview)

	if err != nil {
		returnErrorResponseDebug(c, errors.InternalError, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, uploadUrl)
}

func returnErrorResponseDebug(c *gin.Context, httpError *errors.HttpError, err error) {
	log.Printf("Failed to perform an action, %s\n", err.Error())
	returnErrorResponse(c, httpError)
}

func returnErrorResponse(c *gin.Context, httpError *errors.HttpError) {
	c.IndentedJSON(
		httpError.Status(),
		httpError.Error())
}

func parseUint(id string) (uint64, error) {
	return strconv.ParseUint(id, 10, 32)
}
