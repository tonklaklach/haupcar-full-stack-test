package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"

	"haupcar/backend/internal/handler"
	"haupcar/backend/internal/repository"
	"haupcar/backend/internal/service"
)

func newApp(db *gorm.DB) (*echo.Echo, error) {
	repo := repository.New(db)
	if err := repo.Migrate(); err != nil {
		return nil, err
	}
	cars := handler.NewCar(service.NewCar(repo))

	app := echo.New()
	app.Use(middleware.Recover(), middleware.RequestLogger(), middleware.BodyLimit(16<<10))
	app.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	app.GET("/api/cars", cars.List)
	app.GET("/api/cars/:id", cars.Get)
	app.POST("/api/cars", cars.Create)
	app.PUT("/api/cars/:id", cars.Update)
	app.DELETE("/api/cars/:id", cars.Delete)
	return app, nil
}
