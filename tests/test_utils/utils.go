package test_utils

import (
	"example/wspinapp-backend/pkg/common/schema"
	"gorm.io/gorm"
	"time"
)

var Now = time.Date(2023, 1, 1, 10, 0, 0, 0, time.Local)

var Wall = schema.Wall{
	Model: gorm.Model{
		ID:        1,
		CreatedAt: Now,
		UpdatedAt: Now},
	Holds: []schema.Hold{{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: Now,
			UpdatedAt: Now},
		WallID: 1,
		X:      120.03,
		Y:      256.43,
		Shape:  "Circle",
	}, {
		Model: gorm.Model{
			ID:        2,
			CreatedAt: Now,
			UpdatedAt: Now},
		WallID: 1,
		X:      120.03,
		Y:      36.43,
	}},
	ImageUrl:        "",
	ImagePreviewUrl: "",
}

var Wall2 = schema.Wall{
	Model: gorm.Model{
		ID:        2,
		CreatedAt: Now,
		UpdatedAt: Now},
	Holds: []schema.Hold{{
		Model: gorm.Model{
			ID:        3,
			CreatedAt: Now,
			UpdatedAt: Now},
		WallID: 2,
	}},
}
