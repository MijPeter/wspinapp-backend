package common

import (
	"gorm.io/gorm"
)

type Hold struct {
	gorm.Model
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type Wall struct {
	gorm.Model
	Holds []Hold `json:"holds" gorm:"many2many:wall_holds;"`
	Image string `json:"image"`
}

type Route struct {
	gorm.Model
	Holds  []Hold `json:"holds" gorm:"many2many:route_holds;"`
	WallID uint   `json:"wall"` // probably should be wallId here instead
}
