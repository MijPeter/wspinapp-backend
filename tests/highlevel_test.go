package tests

import (
	"bytes"
	"encoding/json"
	"example/wspinapp-backend/pkg/common"
	"example/wspinapp-backend/pkg/common/schema"
	"example/wspinapp-backend/pkg/controller"
	"example/wspinapp-backend/pkg/services"
	"example/wspinapp-backend/tests/test_utils"
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

func resetDb() {
	clearDb()
	db.AutoMigrate(
		&schema.Wall{},
		&schema.Route{},
		&schema.Hold{},
	)
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
func clearDb() {
	tables, _ := db.Migrator().GetTables()
	for _, table := range tables {
		db.Migrator().DropTable(table)
	}
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

	clearDb()
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

func TestAddWall(t *testing.T) {
	resetDb()
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
	addWallJson, _ := json.Marshal(newWall)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/walls", bytes.NewReader(addWallJson))
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name(), w.Body.String()), w.Body.String())
}

func TestGetWall(t *testing.T) {
	resetDb()
	db.Create(&test_utils.Wall)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/walls/1", nil)
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name(), w.Body.String()), w.Body.String())
}

func TestUpdateWall(t *testing.T) {
	resetDb()
	db.Create(&test_utils.Wall)
	db.Create(&test_utils.Wall2)

	updatedWall := schema.Wall{
		Model: gorm.Model{
			ID:        10,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now()},
		Holds: []schema.Hold{{
			Model: gorm.Model{
				ID:        1,
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
			X:      120.03,
			Y:      36.43,
		}},
		ImageUrl:        "abcd",
		ImagePreviewUrl: "efgh",
	}

	updateWallJson, _ := json.Marshal(&updatedWall)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/walls/1", bytes.NewReader(updateWallJson))
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name()+"_update", w.Body.String()), w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/walls/1", nil)
	req.SetBasicAuth("wspinapp", "wspinapp")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, test_utils.Golden(t.Name()+"_get", w.Body.String()), w.Body.String())
}
