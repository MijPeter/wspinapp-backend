package tests

import (
	"bytes"
	"encoding/json"
	"example/wspinapp-backend/pkg/common"
	"example/wspinapp-backend/pkg/common/adapters/imgrepository"
	"example/wspinapp-backend/pkg/common/schema"
	"example/wspinapp-backend/pkg/controller"
	"example/wspinapp-backend/pkg/services"
	"example/wspinapp-backend/tests/test_utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type ImageRepositoryMock struct {
	returnString string
}

func (i *ImageRepositoryMock) Upload(_ interface{}) (string, error) {
	return i.returnString, nil
}

func (i *ImageRepositoryMock) Assets() ([]imgrepository.Asset, error) {
	return []imgrepository.Asset{}, nil
}

func (i *ImageRepositoryMock) DeleteAssets(_ []string) error {
	return nil
}

var router *gin.Engine
var db *gorm.DB

func initDb() {
	db = common.InitDbWithConfig(&gorm.Config{
		NowFunc: func() time.Time {
			return test_utils.Now
		},
	})
}

/*
TODO think about it if it's possible to reset autoincrement count for tables, currently even though entities are deleted, new entities get nextVal IDs

	so instead of getting ID = 1, we get ID = 6 because we've created 5 entities in previous tests. this should be fixed
	but it's not thaaaat troublesome
	maybe instead of using autoincrement autogenerate IDs (uuid ids) and use random generator with set seed in tests, that should be easy to do but maybe there is something better:D
*/
func clearDb() {
	db.Exec("DELETE FROM route_holds WHERE 1=1;")
	db.Exec("DELETE FROM route_start_holds WHERE 1=1;")
	db.Exec("DELETE FROM route_top_hold WHERE 1=1;")

	db.Unscoped().Where("1 = 1").Delete(schema.Route{})
	db.Unscoped().Where("1 = 1").Delete(schema.Hold{})
	db.Unscoped().Where("1 = 1").Delete(schema.Wall{})
}

// TODO move all that util setup methods to some other file so it's easier to read testss
func initEnv() {
	os.Setenv("POSTGRES_USER", "test_user")
	os.Setenv("POSTGRES_PASSWORD", "test_pass")
	os.Setenv("POSTGRES_DB", "test_db")
	if os.Getenv("POSTGRES_HOST") == "" {
		os.Setenv("POSTGRES_HOST", "0.0.0.0")
	}
	os.Setenv("POSTGRES_PORT", "5432")
}

func TestMain(m *testing.M) {
	log.Println("Setting up tests")
	initEnv()
	initDb()

	router = gin.Default()
	router.SetTrustedProxies(nil)
	router.MaxMultipartMemory = 8 << 20 // 8MiB
	service := services.New(db, &ImageRepositoryMock{"https://static.wikia.nocookie.net"})
	controller.RegisterRoutes(router, service.WebService)

	log.Println("Running tests")

	m.Run()
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestWallAuth(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/walls", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "", w.Body.String())
}

func TestAddWallTwice(t *testing.T) {
	clearDb()
	addWallJson, _ := json.Marshal(&test_utils.WallWithHolds)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/walls", bytes.NewReader(addWallJson))
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name(), w.Body.String()), w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/walls", bytes.NewReader(addWallJson))
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name()+"_2", w.Body.String()), w.Body.String())
}

func TestAddWall(t *testing.T) {
	clearDb()
	newWall := schema.Wall{
		Model: gorm.Model{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now()},
		Holds: []schema.Hold{{ // TODO fill all fields of one of holds to verify if it's saved correctly
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
	addWallJson, _ := json.Marshal(&newWall)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/walls", bytes.NewReader(addWallJson))
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name(), w.Body.String()), w.Body.String())
}

func TestGetWall(t *testing.T) {
	clearDb()
	db.Create(&test_utils.WallWithHolds)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/walls/%d", test_utils.WallWithHolds.ID), nil)
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name(), w.Body.String()), w.Body.String())
}

func TestUpdateWall(t *testing.T) {
	clearDb()
	db.Create(&test_utils.WallFull)
	db.Create(&test_utils.WallWithHolds)

	updatedWall := schema.Wall{
		Model: gorm.Model{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now()},
		Holds: []schema.Hold{{
			Model: gorm.Model{
				ID:        9,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now()},
			WallID: 123,
			X:      300.00,
			Y:      250.43,
			Shape:  "Circle",
		}, {
			Model: gorm.Model{
				ID:        3,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now()},
			WallID: 145,
			X:      777.03,
			Y:      36.43,
		}},
		ImageUrl:        "abcd",
		ImagePreviewUrl: "efgh",
	}

	updateWallJson, _ := json.Marshal(&updatedWall)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/walls/%d", test_utils.WallFull.ID), bytes.NewReader(updateWallJson))
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name(), w.Body.String()), w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/walls/%d", test_utils.WallFull.ID), nil)
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name(), w.Body.String()), w.Body.String())
}

func TestAddRoute(t *testing.T) {
	clearDb()
	db.Create(&test_utils.WallManyHolds)
	holds := test_utils.WallManyHolds.Holds

	assert.Equal(t, 6, len(holds))

	route := schema.Route{
		Holds:      holds[0:4],
		StartHolds: holds[1:3],
		TopHold:    holds[0:1],
	}

	routeJson, _ := json.Marshal(&route)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/walls/%d/routes", test_utils.WallManyHolds.ID), bytes.NewReader(routeJson))
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name(), w.Body.String()), w.Body.String())
}

func TestAddRouteNoHolds(t *testing.T) {
	clearDb()
	db.Create(&test_utils.WallManyHolds)
	holds := test_utils.WallManyHolds.Holds

	assert.Equal(t, 6, len(holds))

	route := schema.Route{}

	routeJson, _ := json.Marshal(&route)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", fmt.Sprintf("/walls/%d/routes", test_utils.WallManyHolds.ID), bytes.NewReader(routeJson))
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name(), w.Body.String()), w.Body.String())
}

func TestUpdateWallDeletesRoutesForHoldsThatAreDeleted(t *testing.T) {
	clearDb()
	db.Create(&test_utils.WallManyHolds)
	holds := test_utils.WallManyHolds.Holds
	assert.Equal(t, 6, len(holds))

	route := schema.Route{
		Holds:      holds[0:4],
		StartHolds: holds[1:3],
		TopHold:    holds[0:1],
		WallID:     test_utils.WallManyHolds.ID,
	}
	db.Create(&route)

	routeThatWillBeDeleted := schema.Route{
		Holds:      holds[3:6],
		StartHolds: holds[3:5],
		TopHold:    holds[5:6],
		WallID:     test_utils.WallManyHolds.ID,
	}
	db.Create(&routeThatWillBeDeleted)

	updatedWall := test_utils.WallManyHolds
	updatedWall.Holds = holds[0:5]

	updateWallJson, _ := json.Marshal(&updatedWall)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/walls/%d", test_utils.WallManyHolds.ID), bytes.NewReader(updateWallJson))
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name()+"_wall", w.Body.String()), w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/walls/%d/routes", test_utils.WallManyHolds.ID), nil)
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name()+"_routes", w.Body.String()), w.Body.String())

	// assert that intermediary tables still contain associations (route is soft deleted)
	var routeHolds, routeStartHolds, routeTopHold []uint
	db.Select("hold_id").Table("route_holds").Where("route_id = ?", routeThatWillBeDeleted.ID).Find(&routeHolds)
	db.Select("hold_id").Table("route_start_holds").Where("route_id = ?", routeThatWillBeDeleted.ID).Find(&routeStartHolds)
	db.Select("hold_id").Table("route_top_hold").Where("route_id = ?", routeThatWillBeDeleted.ID).Find(&routeTopHold)

	assert.True(t, len(routeHolds) == len(routeThatWillBeDeleted.Holds))
	assert.True(t, len(routeStartHolds) == len(routeThatWillBeDeleted.StartHolds))
	assert.True(t, len(routeTopHold) == len(routeThatWillBeDeleted.TopHold))

	// assert route is soft deleted
	assert.Equal(t, nil, db.Find(&schema.Route{}, routeThatWillBeDeleted.ID).Error)

	// assert removed hold is soft deleted
	assert.Equal(t, nil, db.Find(&schema.Hold{}, holds[5]).Error)
}

func TestGetRoutes(t *testing.T) {
	clearDb()
	db.Create(&test_utils.WallManyHolds)
	holds := test_utils.WallManyHolds.Holds
	assert.Equal(t, 6, len(holds))

	route := schema.Route{
		Holds:      holds[0:4],
		StartHolds: holds[1:3],
		TopHold:    holds[0:1],
		WallID:     test_utils.WallManyHolds.ID,
	}
	db.Create(&route)

	routeOther := schema.Route{
		Holds:      holds[3:6],
		StartHolds: holds[3:5],
		TopHold:    holds[5:6],
		WallID:     test_utils.WallManyHolds.ID,
	}
	db.Create(&routeOther)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/walls/%d/routes", test_utils.WallManyHolds.ID), nil)
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name(), w.Body.String()), w.Body.String())
}

func TestDeleteWall(t *testing.T) {
	clearDb()
	db.Create(&test_utils.WallManyHolds)
	holds := test_utils.WallManyHolds.Holds
	assert.Equal(t, 6, len(holds))

	route := schema.Route{
		Holds:      holds[0:4],
		StartHolds: holds[1:3],
		TopHold:    holds[0:1],
		WallID:     test_utils.WallManyHolds.ID,
	}
	db.Create(&route)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/walls/%d", test_utils.WallManyHolds.ID), nil)
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name()+"_wall", w.Body.String()), w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/walls/%d/routes", test_utils.WallManyHolds.ID), nil)
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name()+"_routes", w.Body.String()), w.Body.String())

	var softDeletedHolds []schema.Hold
	db.Where(schema.Hold{WallID: test_utils.WallManyHolds.ID}).Find(&softDeletedHolds)
	assert.Equal(t, []schema.Hold{}, softDeletedHolds)
}

func TestDeleteRoute(t *testing.T) {
	clearDb()
	db.Create(&test_utils.WallManyHolds)
	holds := test_utils.WallManyHolds.Holds
	assert.Equal(t, 6, len(holds))

	route := schema.Route{
		Holds:      holds[0:4],
		StartHolds: holds[1:3],
		TopHold:    holds[0:1],
		WallID:     test_utils.WallManyHolds.ID,
	}
	db.Create(&route)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/walls/%d/routes/%d", test_utils.WallManyHolds.ID, route.ID), nil)
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name(), w.Body.String()), w.Body.String())
}

func TestGetRoute(t *testing.T) {
	clearDb()
	db.Create(&test_utils.WallManyHolds)
	holds := test_utils.WallManyHolds.Holds
	assert.Equal(t, 6, len(holds))

	route := schema.Route{
		Holds:      holds[0:4],
		StartHolds: holds[1:3],
		TopHold:    holds[0:1],
		WallID:     test_utils.WallManyHolds.ID,
	}
	db.Create(&route)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/walls/%d/routes/%d", test_utils.WallManyHolds.ID, route.ID), nil)
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name(), w.Body.String()), w.Body.String())
}

func TestUpdateRoute(t *testing.T) {
	clearDb()
	db.Create(&test_utils.WallManyHolds)
	holds := test_utils.WallManyHolds.Holds
	assert.Equal(t, 6, len(holds))

	route := schema.Route{
		Holds:      holds[0:4],
		StartHolds: holds[1:3],
		TopHold:    holds[0:1],
		WallID:     test_utils.WallManyHolds.ID,
	}
	db.Create(&route)

	routeJson, _ := json.MarshalIndent(&route, "", "\t")
	routeJsonString := string(routeJson)
	assert.Equal(t, test_utils.Golden(t.Name()+"_before_update", routeJsonString), routeJsonString)

	routeUpdated := schema.Route{
		Holds:      holds[3:6],
		StartHolds: holds[3:5],
		TopHold:    holds[5:6],
		WallID:     test_utils.WallManyHolds.ID,
	}

	updatedRouteJson, _ := json.Marshal(&routeUpdated)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/walls/%d/routes/%d", test_utils.WallManyHolds.ID, route.ID), bytes.NewReader(updatedRouteJson))
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name()+"_after_update", w.Body.String()), w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/walls/%d/routes/%d", test_utils.WallManyHolds.ID, route.ID), nil)
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name()+"_after_update", w.Body.String()), w.Body.String())
}
