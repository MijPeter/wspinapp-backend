package test_utils

import (
	"example/wspinapp-backend/pkg/common/schema"
	"time"
)

var Now, _ = time.Parse(time.RFC3339, "2006-01-02T15:04:05+01:00")

var WallMinimal = schema.Wall{}

var WallWithHolds = schema.Wall{
	Holds: []schema.Hold{{
		X:     120.03,
		Y:     256.43,
		Shape: "Circle",
	}, {
		X: 120.03,
		Y: 36.43,
	}},
}

var WallFull = schema.Wall{
	Holds: []schema.Hold{{
		X:     666.03,
		Y:     21.37,
		Shape: "Circle",
	}, {
		X: 120.03,
		Y: 36.43,
	}},
	ImageUrl:        "abcde",
	ImagePreviewUrl: "fghij",
}

var WallManyHolds = schema.Wall{
	Holds: []schema.Hold{{
		X:     13.03,
		Y:     1.43,
		Shape: "Circle",
	}, {
		X: 3,
		Y: 36.43,
	}, {
		X:     13.03,
		Y:     1.43,
		Shape: "Circle",
	}, {
		X: 3,
		Y: 36.43,
	}, {
		X:     13.03,
		Y:     1.43,
		Shape: "Circle",
	}, {
		X: 3,
		Y: 36.43,
	}},
}
