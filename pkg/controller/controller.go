package controller

import (
	"example/wspinapp-backend/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type routesHandler struct {
	service  services.WebService
	validate *validator.Validate
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func RegisterRoutes(r *gin.Engine, service services.WebService) {
	h := &routesHandler{
		service:  service,
		validate: validator.New(),
	}

	publicRouter := r.Group("")
	publicRouter.GET("/ping", h.Pong)

	publicRouter.GET("/metrics", prometheusHandler())

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
