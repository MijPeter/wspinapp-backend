package controller

import (
	"example/wspinapp-backend/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type routesHandler struct {
	service  services.Service
	validate *validator.Validate
}

func RegisterRoutes(r *gin.Engine, service services.Service) {
	h := &routesHandler{
		service:  service,
		validate: validator.New(),
	}

	publicRouter := r.Group("")
	publicRouter.GET("/ping", h.Pong)

	router := r.Group("/walls")
	router.Use(gin.BasicAuth(gin.Accounts{
		"wspinapp": "wspinapp",
	}))

	router.POST("", h.AddWall)
	router.GET("", h.GetWalls)
	router.GET("/:wallId", h.GetWall)
	router.PUT("/:wallId", h.UpdateWall)
	router.DELETE("/:wallId", h.DeleteWall)
	router.GET("/:wallId/routes", h.GetRoutes)
	router.POST("/:wallId/routes", h.AddRoute)
	router.GET("/:wallId/routes/:routeId", h.GetRoute)
	router.PUT("/:wallId/routes/:routeId", h.UpdateRoute)
	router.DELETE("/:wallId/routes/:routeId", h.DeleteRoute)
	router.PATCH("/:wallId/image", h.UploadImageFull)
	router.PATCH("/:wallId/imagepreview", h.UploadImagePreview)
}
