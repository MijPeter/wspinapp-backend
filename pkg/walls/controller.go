package walls

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type wallRoutesHandler struct {
	database *gorm.DB
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	h := &wallRoutesHandler{
		database: db,
	}
	router := r.Group("/walls")

	router.POST("/", h.AddWall)
	router.GET("/", h.GetWalls)
	router.GET("/:wallId", h.GetWall)
	router.GET("/:wallId/routes", h.GetRoutes)
}
