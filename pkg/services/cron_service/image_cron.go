package cron_service

import (
	"example/wspinapp-backend/pkg/common/adapters/imgrepository"
	"gorm.io/gorm"
	"log"
	"time"
)

const HOURS = 24 // Run is activated every HOURS

type CronService struct {
	database        *gorm.DB
	imageRepository imgrepository.ImageRepository
}

func New(db *gorm.DB, imageRepository imgrepository.ImageRepository) CronService {
	return CronService{database: db, imageRepository: imageRepository}
}

func (c *CronService) Run() {
	go c.runImageCron()
	go c.runWallCron()
}

func (c *CronService) runImageCron() {
	// in theory should be run at 3 am, but this is not that important right now
	for true {
		time.Sleep(HOURS * time.Hour)
		log.Println("DING DING DING, running deleteAllDanglingImages cron job")
		c.deleteAllDanglingImages()
	}
}

func (c *CronService) deleteAllDanglingImages() {
	var activeImageUrls, imageUrls, imagePreviewUrls []string
	c.database.Select("image_url").Table("walls").Find(&imageUrls)
	c.database.Select("image_preview_url").Table("walls").Find(&imagePreviewUrls)

	activeImageUrls = append(imageUrls, imagePreviewUrls...)

	log.Println("Active images that won't be deleted:")

	for _, img := range activeImageUrls {
		log.Println(img)
	}

	imgRepositoryAssets, err := c.imageRepository.Assets()
	if err != nil {
		log.Printf("Couldn't list assets. Reason\n%s", err.Error()) // some more serious log here
		return
	}
	log.Println("Printing all assets:")
	for _, asset := range imgRepositoryAssets {
		log.Println(asset)
	}

	activeImageUrlsMap := make(map[string]struct{})
	for _, activeImageUrl := range activeImageUrls {
		activeImageUrlsMap[activeImageUrl] = struct{}{}
	}

	var toBeDeletedAssetsID []string
	for _, asset := range imgRepositoryAssets {
		_, exists := activeImageUrlsMap[asset.Url]
		if !exists {
			toBeDeletedAssetsID = append(toBeDeletedAssetsID, asset.ID)
		}
	}

	log.Printf("About to delete %d assets", len(toBeDeletedAssetsID))

	err = c.imageRepository.DeleteAssets(toBeDeletedAssetsID)
	if err != nil {
		log.Printf("Failed to delete assets: Reason\n%s", err.Error())
	} else {
		log.Println("Successfuly deleted assets :)")
	}
}
