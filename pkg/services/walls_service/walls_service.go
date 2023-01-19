package walls_service

import (
	"example/wspinapp-backend/pkg/common/errors"
	"example/wspinapp-backend/pkg/common/schema"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
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
func GetRoutes(db *gorm.DB, wallId uint) []schema.Route {
	var wallRoutes []schema.Route
	db.Preload(clause.Associations).Where(schema.Route{WallID: wallId}).Find(&wallRoutes)
	return wallRoutes
}

func AddRoute(db *gorm.DB, route *schema.Route, wallId uint) error {
	var wall schema.Wall
	err := db.Preload(clause.Associations).First(&wall, wallId).Error

	if err != nil {
		return err
	}

	holdsMap := make(map[uint]schema.Hold)

	for _, hold := range wall.Holds {
		holdsMap[hold.ID] = hold
	}

	var holds []schema.Hold
	for _, hold := range route.Holds {
		realHold, ok := holdsMap[hold.ID]
		if !ok {
			return errors.New("Hold doesn't belong to wall", 400)
		} else {
			holds = append(holds, realHold)
		}
	}

	var startHolds []schema.Hold
	for _, hold := range route.StartHolds {
		realHold, ok := holdsMap[hold.ID]
		if !ok {
			return errors.New("Hold doesn't belong to wall", 400)
		} else {
			startHolds = append(startHolds, realHold)
		}
	}

	if len(route.TopHold) > 1 {
		return errors.New("Too many top holds", 400)
	} else if len(route.TopHold) == 1 {
		realHold, ok := holdsMap[route.TopHold[0].ID]
		if !ok {
			return errors.New("Hold doesn't belong to wall", 400)
		} else {
			route.TopHold[0] = realHold
		}
	}

	route.Holds = holds
	route.StartHolds = startHolds
	route.WallID = wallId
	log.Println(len(route.Holds))
	return db.Create(route).Error
}
