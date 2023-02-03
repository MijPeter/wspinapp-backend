package services

import (
	"example/wspinapp-backend/pkg/common/adapters/imgrepository"
	"example/wspinapp-backend/pkg/services/cron_service"
	"example/wspinapp-backend/pkg/services/walls_service"
	"gorm.io/gorm"
)

type WebService struct {
	WallsService walls_service.WallsService
}

type CronService struct {
	CronService cron_service.CronService
}

func New(db *gorm.DB, imgRepository imgrepository.ImageRepository) (WebService, CronService) {
	return WebService{
			WallsService: walls_service.New(db, imgRepository),
		},
		CronService{
			CronService: cron_service.New(db, imgRepository),
		}
}
