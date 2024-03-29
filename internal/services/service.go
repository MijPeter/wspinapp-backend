package services

import (
	"example/wspinapp-backend/internal/common/adapters/imgrepository"
	"example/wspinapp-backend/internal/services/cron_service"
	"example/wspinapp-backend/internal/services/walls_service"
	"gorm.io/gorm"
)

type Service struct {
	WebService  WebService
	CronService cron_service.CronService
}

type WebService struct {
	WallsService walls_service.WallsService
}

func New(db *gorm.DB, imgRepository imgrepository.ImageRepository) Service {
	return Service{
		WebService:  WebService{WallsService: walls_service.New(db, imgRepository)},
		CronService: cron_service.New(db, imgRepository),
	}
}
