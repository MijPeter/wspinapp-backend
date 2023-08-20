package controller

import (
	"example/wspinapp-backend/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
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

func attachMetrics() gin.HandlerFunc {
	buckets := []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

	responseTimeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "wspinapp",
		Name:      "http_server_request_duration_seconds",
		Help:      "Histogram of response time for handler in seconds",
		Buckets:   buckets,
	}, []string{"route", "method", "status_code"})

	prometheus.MustRegister(responseTimeHistogram)

	return func(c *gin.Context) {
		start := time.Now()

		// Let the request go through the handler chain
		c.Next()

		// Now calculate the metrics after the handler has finished processing
		duration := time.Since(start)

		// Get route (path), method and status code from the current context
		route := c.Request.URL.Path // Or use c.FullPath() for the matched route pattern
		method := c.Request.Method
		statusCode := strconv.Itoa(c.Writer.Status())

		responseTimeHistogram.WithLabelValues(route, method, statusCode).Observe(duration.Seconds())
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
	}), attachMetrics())

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
