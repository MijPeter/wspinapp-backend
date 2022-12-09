package common

type Position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}

type Hold struct {
	Id       string   `json:"id"`
	Position Position `json:"position"`
}

type Wall struct {
	Id    string `json:"id"`
	Holds []Hold `json:"holds"`
	Image string `json:"image"`
}

type Route struct {
	Id     string `json:"id"`
	Holds  []Hold `json:"holds"`
	WallId string `json:"wall"` // probably should be wallId here instead
}
