package cron_service

import (
	"example/wspinapp-backend/pkg/common/schema"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

const WALL_HOURS = 156

func (c *CronService) runWallCron() {
	// in theory should be run at 3 am, but this is not that important right now
	for true {
		log.Println("DING DING DING, running deleteAllDanglingWalls cron job")
		c.deleteAllDanglingWalls()
		time.Sleep(WALL_HOURS * time.Hour)
	}
}

func (c *CronService) deleteAllDanglingWalls() {
	// Get all walls without imageUrls that are older than 24 hours
	err := c.database.Preload(clause.Associations).Where("image_url = ?", "").Delete(&schema.Wall{}).Error
	if err != nil {
		log.Printf("Failed to delete dangling walls. Reason\n%s", err.Error())
		return
	}
	log.Println("Successfuly deleted dangling walls")
}
