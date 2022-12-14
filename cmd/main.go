package main

import (
	"example/wspinapp-backend/pkg/common"
	"example/wspinapp-backend/pkg/common/adapters/imgrepository"
	"example/wspinapp-backend/pkg/controller"
	"github.com/gin-gonic/gin"
)

var basicAuth = gin.BasicAuth(gin.Accounts{
	"wspinapp": "wspinapp",
})

func main() {
	db := common.ConnectDb()
	db.AutoMigrate(
		&common.Hold{},
		&common.Wall{},
		&common.Route{},
	)
	// todo probably create some simple adapter for db for cleanliness sake

	router := gin.Default()
	router.Use(basicAuth)               // TODO add google account auth
	router.MaxMultipartMemory = 8 << 20 // 8MiB

	controller.RegisterRoutes(router, db, imgrepository.New())
	router.Run()
}
