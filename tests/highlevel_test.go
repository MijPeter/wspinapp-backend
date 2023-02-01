package tests

import (
	"bytes"
	"encoding/json"
	"example/wspinapp-backend/pkg/common"
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

var router *gin.Engine
var db *gorm.DB

func initDb() {
	db = common.InitDbWithConfig(&gorm.Config{
		NowFunc: func() time.Time {
			return test_utils.Now
		},
	})
}

func clearDb() {
	db.Unscoped().Where("1 = 1").Delete(schema.Hold{})
	db.Unscoped().Where("1 = 1").Delete(schema.Route{})
	db.Unscoped().Where("1 = 1").Delete(schema.Wall{})
}

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
	controller.RegisterRoutes(router, service)

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
	assert.Equal(t, test_utils.Golden(t.Name()+"_1", w.Body.String()), w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/walls/%d", test_utils.WallFull.ID), nil)
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name()+"_2", w.Body.String()), w.Body.String())
}
