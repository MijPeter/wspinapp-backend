package schema

import (
	"gorm.io/gorm"
	"mime/multipart"
)

type Hold struct {
	gorm.Model
	X      float32 `json:"X" gorm:"not null"`
	Y      float32 `json:"Y" gorm:"not null"`
	Size   float32 `json:"Size" gorm:"not null"`
	WallID uint    `json:"WallID" gorm:"not null"`
	Shape  string  `json:"Shape" gorm:"not null;default:circle"`
	Angle  float32 `json:"Angle"`
}

type Wall struct {
	gorm.Model
	Holds    []Hold `json:"Holds"`
	ImageUrl string `json:"ImageUrl"`
}

type Route struct {
	gorm.Model
	Holds      []Hold `json:"Holds" gorm:"many2many:route_holds"`
	StartHolds []Hold `json:"StartHolds" gorm:"many2many:route_holds"`
	TopHold    []Hold `json:"TopHold" gorm:"many2many:route_holds"`
	WallID     uint   `json:"WallID" gorm:"not_null"`
}

type File struct {
	File multipart.File `json:"file,omitempty" validate:"required"`
}

type Form struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type Url struct {
	Url string `json:"url,omitempty" validate:"required"`
}
