package service

import (
	"github.com/gin-gonic/gin"

	"sample-project/storage"
)

type (
	APIv1 struct {
		store       storage.Storager
		redisClient storage.RedisStorager
	}
	Response struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
)

func Start(store storage.Storager, redisClient storage.RedisStorager) error {
	api := APIv1{store: store, redisClient: redisClient}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/login", api.login)
	router.POST("/logout", api.tokenAuthMiddleware(), api.logout)

	router.POST("/drones", api.tokenAuthMiddleware(), api.createDrone)
	router.GET("/drones", api.tokenAuthMiddleware(), api.getAllDrones)
	return router.Run(":8080")
}
