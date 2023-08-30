package controller

import (
	"example/wspinapp-backend/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
)

type routesHandler struct {
	service             services.WebService
	validate            *validator.Validate
	metricsHandler      http.Handler
	metricsHistogramVec *prometheus.HistogramVec
}

func prometheusHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func initPrometheus() *prometheus.HistogramVec {
	responseTimeHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "wspinapp",
		Name:      "http_server_request_duration_seconds",
		Help:      "Histogram of response time for handler in seconds",
	}, []string{"route", "method", "status_code"})

	prometheus.MustRegister(responseTimeHistogram)
	return responseTimeHistogram
}

func attachMetrics(responseTimeHistogram *prometheus.HistogramVec) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Let the request go through the handler chain
		c.Next()

		// Now calculate the metrics after the handler has finished processing
		duration := time.Since(start)

		// Get route (path), method and status code from the current context
		route := c.FullPath()
		method := c.Request.Method
		statusCode := strconv.Itoa(c.Writer.Status())

		responseTimeHistogram.WithLabelValues(route, method, statusCode).Observe(duration.Seconds())
	}
}

func RegisterRoutes(r *gin.Engine, service services.WebService) {
	h := &routesHandler{
		service:             service,
		validate:            validator.New(),
		metricsHandler:      promhttp.Handler(),
		metricsHistogramVec: initPrometheus(),
	}

	publicRouter := r.Group("")
	publicRouter.GET("/ping", h.Pong)

	publicRouter.GET("/metrics", prometheusHandler(h.metricsHandler))

	router := r.Group("/walls")
	router.Use(gin.BasicAuth(gin.Accounts{
		"wspinapp": "wspinapp",
	}), attachMetrics(h.metricsHistogramVec))

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
