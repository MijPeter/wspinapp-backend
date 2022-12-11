package main

import (
	"example/wspinapp-backend/pkg/common"
	"example/wspinapp-backend/pkg/common/adapters/imgrepository"
	"example/wspinapp-backend/pkg/images"
	"example/wspinapp-backend/pkg/walls"
	"github.com/gin-gonic/gin"
)

var basicAuth = gin.BasicAuth(gin.Accounts{
	"wspinapp": "wspinapp",
})

func main() {
	common.LoadEnvironmentVariables()

	db := common.ConnectDb()
	db.AutoMigrate(
		&common.Hold{},
		&common.Wall{},
		&common.Route{},
	)
	// todo probably create some simple adapter for db for cleanliness sake

	router := gin.Default()
	router.Use(basicAuth)

	walls.RegisterRoutes(router, db)
	images.RegisterRoutes(router, db, imgrepository.New())
	router.Run()
}
