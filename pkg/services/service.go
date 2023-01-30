package services

import (
	"example/wspinapp-backend/pkg/common/adapters/imgrepository"
	"example/wspinapp-backend/pkg/services/walls_service"
	"gorm.io/gorm"
)

type Service struct {
	WallsService walls_service.WallsService
}

func New(db *gorm.DB, imgRepository imgrepository.ImageRepository) Service {
	return Service{
		WallsService: walls_service.New(db, imgRepository),
	}
}
