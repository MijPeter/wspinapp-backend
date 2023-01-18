package walls_service

import (
	"example/wspinapp-backend/pkg/common/schema"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO add db persistence service layer (a layer that communicates with db)
// TODO but maybe for now its not needed

func AddWall(db *gorm.DB, wall *schema.Wall) {
	db.Create(wall)
}

func GetWall(db *gorm.DB, wallId uint) (schema.Wall, error) {
	var wall schema.Wall

	err := db.Preload(clause.Associations).First(&wall, wallId).Error
	return wall, err
}

func GetWalls(db *gorm.DB) []schema.Wall {
	var walls []schema.Wall
	db.Preload(clause.Associations).Find(&walls)
	return walls
}

// TODO routes aren't implemented yet
//func GetRoutes(db *gorm.DB, wallId uint) []schema.Route {
//	var wallRoutes []schema.Route
//	db.Preload(clause.Associations).Where(schema.Route{WallID: wallId}).Find(&wallRoutes)
//	return wallRoutes
//}
