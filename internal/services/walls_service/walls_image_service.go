package walls_service

import (
	"example/wspinapp-backend/internal/common/schema"
	"log"
)

func (s *WallsService) UploadFileAndSaveUrlToDb(
	wallId uint,
	file schema.File,
	isPreview bool) (string, error) {

	uploadUrl, err := s.fileUpload(file)

	if err != nil {
		log.Printf("Failed to upload image, %s\n", err.Error())
		return "", err
	}

	err = s.uploadUrlToDb(wallId, uploadUrl, isPreview)

	return uploadUrl, err
}

func (s *WallsService) uploadUrlToDb(wallId uint, url string, isPreview bool) error {
	var wall schema.Wall

	err := s.database.First(&wall, wallId).Error

	if err != nil {
		return err
	}

	if isPreview {
		wall.ImagePreviewUrl = url
	} else {
		wall.ImageUrl = url
	}
	return s.database.Save(&wall).Error
}

func (s *WallsService) fileUpload(file schema.File) (string, error) {
	//upload
	uploadUrl, err := s.imageRepository.Upload(file.File)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}

//https://dev.to/hackmamba/robust-media-upload-with-golang-and-cloudinary-gin-gonic-version-54ii
