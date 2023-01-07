package walls_service

import (
	"example/wspinapp-backend/pkg/common"
	"gorm.io/gorm"
)

// TODO add db persistence service layer (a layer that communicates with db)

func AddWall(db *gorm.DB, wall *common.Wall) {
	db.Create(wall)
}

func GetWall(db *gorm.DB, wallId uint) (common.Wall, error) {
	var wall common.Wall

	err := db.Preload("Holds").First(&wall, wallId).Error
	return wall, err
}

func GetWalls(db *gorm.DB) []common.Wall {
	var walls []common.Wall
	db.Preload("Holds").Find(&walls)
	return walls
}

func GetRoutes(db *gorm.DB, wallId uint) []common.Route {
	var wallRoutes []common.Route
	db.Where(common.Route{WallID: wallId}).Find(&wallRoutes)
	return wallRoutes
}
