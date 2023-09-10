package main

import (
	"example/wspinapp-backend/internal/common"
	"example/wspinapp-backend/internal/common/adapters/imgrepository"
	"example/wspinapp-backend/internal/controller"
	"example/wspinapp-backend/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// todo probably create some simple adapter for db for cleanliness sake
	db := common.InitDb()

	router := gin.Default()
	router.SetTrustedProxies(nil)       // TODO add google account auth
	router.MaxMultipartMemory = 8 << 20 // 8MiB

	service := services.New(db, imgrepository.New())
	controller.RegisterRoutes(router, service.WebService)
	service.CronService.Run()
	router.Run()
}

// GENERAL TODOS

// TODO add google account auth
// TODO add error wrapping and standardize error messages for responses
// TODO remove image when deleting wall
// TODO only remove routes that have holds that are deleted when updating wall
// TODO
// TODO	TESTSTESTSETSE
// TODO TESTESTSETESTSE
// here is some inspiration for tests https://github.com/gothinkster/golang-gin-realworld-example-app
// also here https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
// todo probably create some simple adapter for db for cleanliness sake
