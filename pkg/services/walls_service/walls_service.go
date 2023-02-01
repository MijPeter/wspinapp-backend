package walls_service

import (
	"example/wspinapp-backend/pkg/common/adapters/imgrepository"
	"example/wspinapp-backend/pkg/common/errors"
	"example/wspinapp-backend/pkg/common/schema"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type WallsService struct {
	database        *gorm.DB
	imageRepository imgrepository.ImageRepository
}

func New(db *gorm.DB, imageRepository imgrepository.ImageRepository) WallsService {
	return WallsService{database: db, imageRepository: imageRepository}
}

func (s *WallsService) AddWall(wall *schema.Wall) {
	wall.Model = gorm.Model{}

	for i := range wall.Holds {
		wall.Holds[i].Model = gorm.Model{}
	}

	wall.ImageUrl = ""
	wall.ImagePreviewUrl = ""

	s.database.Create(wall)
}

func (s *WallsService) GetWall(wallId uint) (schema.Wall, error) {
	var wall schema.Wall

	err := s.database.Preload(clause.Associations).First(&wall, wallId).Error
	return wall, err
}

func (s *WallsService) UpdateWall(wallId uint, newWall *schema.Wall) (schema.Wall, error) {
	var stateWall schema.Wall

	err := s.database.Preload(clause.Associations).First(&stateWall, wallId).Error
	if err != nil {
		return stateWall, err
	}

	// delete all routes for a newWall
	s.database.Delete(&schema.Route{}, "wall_id = ?", wallId)
	// TODO delete only those routes that have their holds removed

	// delete holds that are being removed ass well TODO

	currentHolds := make(map[uint]schema.Hold)

	for i := range stateWall.Holds {
		hold := stateWall.Holds[i]
		currentHolds[hold.ID] = hold
	}

	var newStateHolds []schema.Hold
	for i := range newWall.Holds {
		newHold := newWall.Holds[i]
		stateHold, ok := currentHolds[newHold.ID]

		newHold.WallID = wallId
		if !ok {
			newHold.Model = gorm.Model{}
			newStateHolds = append(newStateHolds, newHold)
		} else {
			copyHoldInto(newHold, &stateHold)
			err = s.database.Updates(&stateHold).Error
			newStateHolds = append(newStateHolds, stateHold)
		}
		if err != nil {
			return *newWall, err
		}
	}

	// Delete holds that aren't present in newHolds
	newHolds := make(map[uint]schema.Hold)
	for i := range newWall.Holds {
		hold := newWall.Holds[i]
		newHolds[hold.ID] = hold
	}

	for i := range stateWall.Holds {
		stateHold := stateWall.Holds[i]
		_, ok := newHolds[stateHold.ID]
		if !ok {
			err = s.database.Delete(&stateHold, "id = ?", stateHold.ID).Error
			log.Printf("Deleting hold with id :%d", stateHold.ID)
			if err != nil {
				return *newWall, err
			}
		}
	}

	stateWall.Holds = newStateHolds
	err = s.database.Save(&stateWall).Error
	return stateWall, err
}

func (s *WallsService) DeleteWall(wallId uint) {
	// delete all routes for a wall
	s.database.Delete(&schema.Route{}, "wall_id = ?", wallId)
	s.database.Delete(&schema.Hold{}, "wall_id = ?", wallId)
	s.database.Delete(&schema.Wall{}, wallId)
}

func (s *WallsService) GetWalls() []schema.Wall {
	var walls []schema.Wall
	s.database.Preload(clause.Associations).Find(&walls)
	return walls
}

func (s *WallsService) GetRoutes(wallId uint) []schema.Route {
	var wallRoutes []schema.Route
	s.database.Preload(clause.Associations).Where(schema.Route{WallID: wallId}).Find(&wallRoutes)
	return wallRoutes
}

// TODO refactor this thing beloww
func (s *WallsService) AddRoute(route *schema.Route, wallId uint) error {
	var wall schema.Wall
	err := s.database.Preload(clause.Associations).First(&wall, wallId).Error

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
	return s.database.Create(route).Error
}

func (s *WallsService) GetRoute(wallId uint, routeId uint) (schema.Route, error) {
	// TODO should we validate wallId against routeId?
	var stateRoute schema.Route

	err := s.database.Preload(clause.Associations).First(&stateRoute, routeId).Error
	return stateRoute, err
}

func (s *WallsService) UpdateRoute(route *schema.Route, wallId uint, routeId uint) error {
	var stateRoute schema.Route

	err := s.database.Preload(clause.Associations).First(&stateRoute, routeId).Error

	if err != nil {
		return err
	}

	var wall schema.Wall
	err = s.database.Preload(clause.Associations).First(&wall, wallId).Error

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
			s.database.Model(&stateRoute).Association("TopHold").Replace([]schema.Hold{realHold})
		}
	} else {
		s.database.Model(&stateRoute).Association("TopHold").Replace([]schema.Hold{})
	}
	stateRoute.WallID = wallId
	log.Println(len(stateRoute.Holds))
	s.database.Model(&stateRoute).Association("Holds").Replace(holds)
	s.database.Model(&stateRoute).Association("StartHolds").Replace(startHolds)

	return s.database.Save(stateRoute).Error
}

func (s *WallsService) DeleteRoute(wallId uint, routeId uint) error {
	var stateRoute schema.Route

	err := s.database.First(&stateRoute, routeId).Error

	if err != nil {
		return err
	}

	var wall schema.Wall
	err = s.database.Preload(clause.Associations).First(&wall, wallId).Error

	if err != nil {
		return err
	}

	if stateRoute.WallID != wallId {
		return errors.BadRequest
	}

	// Also delete assosciatio tables
	return s.database.Select(clause.Associations).Delete(&schema.Route{}, routeId).Error
}

func copyHoldInto(from schema.Hold, to *schema.Hold) {
	to.X = from.X
	to.Y = from.Y
	to.Size = from.Size
	to.Shape = from.Shape
	to.Angle = from.Angle
}
