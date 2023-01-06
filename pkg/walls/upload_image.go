package walls

import (
	"example/wspinapp-backend/pkg/common"
	"example/wspinapp-backend/pkg/common/adapters/imgrepository"
	"example/wspinapp-backend/pkg/common/errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *routesHandler) UploadImage(c *gin.Context) {
	uploadedFile, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("Failed to parse given image, %s\n", err.Error())
		c.IndentedJSON(
			errors.InternalError.Status(),
			errors.InternalError.Error())
		return
	}

	uploadUrl, err := h.FileUpload(common.File{File: uploadedFile}, h.imageRepository)

	if err != nil {
		log.Printf("Failed to upload image, %s\n", err.Error())
		c.IndentedJSON(
			errors.InternalError.Status(),
			errors.InternalError.Error())
		return
	}

	err = uploadUrlToDb(uploadUrl, c, h)

	c.IndentedJSON(http.StatusCreated, uploadUrl)
}

func uploadUrlToDb(url string, c *gin.Context, h *routesHandler) error {
	wallId64, err := strconv.ParseUint(c.Param("wallId"), 10, 32)
	wallId := uint(wallId64)

	if err != nil {
		return err
	}
	var wall common.Wall

	err = h.database.First(&wall, wallId).Error

	if err != nil {
		return err
	}

	wall.Image = url
	return h.database.Save(&wall).Error
}

func (h *routesHandler) FileUpload(file common.File, imageRepository imgrepository.ImageRepository) (string, error) {
	//validate
	err := h.validate.Struct(file)
	if err != nil {
		return "", err
	}

	//upload
	uploadUrl, err := imageRepository.Upload(file.File)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}

//https://dev.to/hackmamba/robust-media-upload-with-golang-and-cloudinary-gin-gonic-version-54ii
