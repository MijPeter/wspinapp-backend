package test_utils

import (
	"example/wspinapp-backend/pkg/common/schema"
	"gorm.io/gorm"
	"time"
)

var Wall = schema.Wall{
	Model: gorm.Model{
		ID:        10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now()},
	Holds: []schema.Hold{{
		Model: gorm.Model{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now()},
		WallID: 123,
		X:      120.03,
		Y:      256.43,
		Shape:  "Circle",
	}, {
		Model: gorm.Model{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now()},
		WallID: 145,
		X:      120.03,
		Y:      36.43,
	}},
	ImageUrl:        "abcd",
	ImagePreviewUrl: "efgh",
}
