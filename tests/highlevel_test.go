package tests

import (
	"example/wspinapp-backend/pkg/common"
	"example/wspinapp-backend/pkg/controller"
	"example/wspinapp-backend/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type ImageRepositoryMock struct {
	returnString string
}

func (i *ImageRepositoryMock) Upload(_ interface{}) (string, error) {
	return i.returnString, nil
}

var router *gin.Engine

func TestMain(m *testing.M) {
	log.Println("Setting up tests")

	os.Setenv("POSTGRES_USER", "test_user")
	os.Setenv("POSTGRES_PASSWORD", "test_pass")
	os.Setenv("POSTGRES_DB", "test_db")
	os.Setenv("POSTGRES_HOST", "0.0.0.0")
	os.Setenv("POSTGRES_PORT", "5432")

	db := common.InitDb()

	router = gin.Default()
	router.SetTrustedProxies(nil)
	router.MaxMultipartMemory = 8 << 20 // 8MiB
	service := services.New(db, &ImageRepositoryMock{"https://static.wikia.nocookie.net"})
	controller.RegisterRoutes(router, service)

	log.Println("Running tests")
	m.Run()

	// clear db
	tables, _ := db.Migrator().GetTables()
	for _, table := range tables {
		db.Migrator().DropTable(table)
	}
}

func TestPingRoute(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
