package controller

import (
	"example/wspinapp-backend/pkg/common/adapters/imgrepository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type routesHandler struct {
	database        *gorm.DB
	imageRepository imgrepository.ImageRepository
	validate        *validator.Validate
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB, imageRepository imgrepository.ImageRepository) {
	h := &routesHandler{
		database:        db,
		imageRepository: imageRepository,
		validate:        validator.New(),
	}
	router := r.Group("/walls")

	router.POST("", h.AddWall)
	router.GET("", h.GetWalls)
	router.GET("/:wallId", h.GetWall)
	router.GET("/:wallId/routes", h.GetRoutes)
	router.PATCH("/:wallId/image", h.UploadImage)
}
