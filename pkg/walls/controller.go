package walls

import (
	"example/wspinapp-backend/pkg/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	h := &handler{
		DB: db,
	}
	router := r.Group("/walls")

	router.POST("/", h.AddWall)
	router.GET("/", h.GetWalls)
	router.GET("/:wallId", h.GetWall)
	router.GET("/:wallId/routes", h.GetRoutes)
}

// this should be db :)
var walls = []common.Wall{
	{Id: "0", Holds: []common.Hold{}, Image: "image_1.png"},
}
var routes []common.Route
