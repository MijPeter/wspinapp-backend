package main

import (
	"example/wspinapp-backend/pkg/common"
	"example/wspinapp-backend/pkg/walls"
	"github.com/gin-gonic/gin"
)

var basicAuth = gin.BasicAuth(gin.Accounts{
	"wspinapp": "wspinapp",
})

func main() {
	db := common.ConnectDb()
	router := gin.Default()
	router.Use(basicAuth)

	walls.RegisterRoutes(router, db)

	router.Run()
}