package common

import (
	"gorm.io/gorm"
	"mime/multipart"
)

type Hold struct {
	gorm.Model
	X      int32 `json:"X"`
	Y      int32 `json:"Y"`
	WallID uint  `json:"WallID"`
}

type Wall struct {
	gorm.Model
	Holds []Hold `json:"Holds"`
	Image string `json:"Image"`
}

type Route struct {
	gorm.Model
	Holds  []Hold `json:"Holds" gorm:"many2many:route_holds;"`
	WallID uint   `json:"WallId"`
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
