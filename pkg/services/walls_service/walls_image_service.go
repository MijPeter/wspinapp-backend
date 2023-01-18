package walls_service

import (
	"example/wspinapp-backend/pkg/common/adapters/imgrepository"
	"example/wspinapp-backend/pkg/common/schema"
	"gorm.io/gorm"
	"log"
)

func UploadFileAndSaveUrlToDb(
	db *gorm.DB,
	imageRepository imgrepository.ImageRepository,
	wallId uint,
	file schema.File) (string, error) {

	uploadUrl, err := FileUpload(file, imageRepository)

	if err != nil {
		log.Printf("Failed to upload image, %s\n", err.Error())
		return "", err
	}

	err = uploadUrlToDb(wallId, uploadUrl, db)

	return uploadUrl, err
}

func uploadUrlToDb(wallId uint, url string, db *gorm.DB) error {
	var wall schema.Wall

	err := db.First(&wall, wallId).Error

	if err != nil {
		return err
	}

	wall.ImageUrl = url
	return db.Save(&wall).Error
}

func FileUpload(file schema.File, imageRepository imgrepository.ImageRepository) (string, error) {
	//upload
	uploadUrl, err := imageRepository.Upload(file.File)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}

//https://dev.to/hackmamba/robust-media-upload-with-golang-and-cloudinary-gin-gonic-version-54ii
