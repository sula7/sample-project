package service

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	g := e.Group("/api/v1")

	g.POST("/login", api.login)
	g.POST("/logout", api.logout, api.tokenAuthMiddleware)

	g.POST("/drones", api.createDrone, api.tokenAuthMiddleware)
	g.GET("/drones", api.getAllDrones, api.tokenAuthMiddleware)
	return e.Start(":8080")
}
