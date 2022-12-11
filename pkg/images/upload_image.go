package images

import (
	"example/wspinapp-backend/pkg/common"
	"example/wspinapp-backend/pkg/common/adapters/imgrepository"
	"example/wspinapp-backend/pkg/common/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (h *imageRoutesHandler) UploadImage(c *gin.Context) {
	uploadedFile, _, err := c.Request.FormFile("file")
	if err != nil {
		c.IndentedJSON(
			errors.InternalError.Status(),
			errors.InternalError.Error())
		return
	}

	uploadUrl, err := FileUpload(common.File{File: uploadedFile}, h.imageRepository)

	if err != nil {
		c.IndentedJSON(
			errors.InternalError.Status(),
			errors.InternalError.Error())
		return
	}
	c.IndentedJSON(http.StatusCreated, uploadUrl)
}

var (
	validate = validator.New()
)

func FileUpload(file common.File, imageRepository imgrepository.ImageRepository) (string, error) {
	//validate
	err := validate.Struct(file)
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
