package images

import (
	"example/wspinapp-backend/pkg/common/adapters/imgrepository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type imageRoutesHandler struct {
	database        *gorm.DB
	imageRepository imgrepository.ImageRepository
	validate        *validator.Validate
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB, imageRepository imgrepository.ImageRepository) {
	routeHandler := &imageRoutesHandler{
		database:        db,
		imageRepository: imageRepository,
		validate:        validator.New(),
	}
	router := r.Group("/images")

	router.POST("/", routeHandler.UploadImage)
	router.GET("/", routeHandler.DownloadImage)
}

// do as described here :))
//https://dev.to/hackmamba/robust-media-upload-with-golang-and-cloudinary-gin-gonic-version-54ii
