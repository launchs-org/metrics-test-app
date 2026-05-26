package main

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"metrics-sample-app/internal/handler"
	"metrics-sample-app/internal/stress"
)

//go:embed static
var staticFiles embed.FS

func main() {
	stressState := &stress.State{}
	stressHandler := handler.NewStressHandler(stressState)

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	staticFS, _ := fs.Sub(staticFiles, "static")
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(staticFS))))

	api := e.Group("/api")
	api.GET("/env", handler.Env)
	api.GET("/workdir", handler.Workdir)
	api.GET("/sysinfo", handler.SysInfo)
	api.GET("/dir", handler.DirExplorer)
	api.GET("/file", handler.FileServe)
	api.GET("/stress/status", stressHandler.Status)
	api.POST("/stress/cpu/start", stressHandler.CPUStart)
	api.POST("/stress/cpu/stop", stressHandler.CPUStop)
	api.POST("/stress/mem/start", stressHandler.MemStart)
	api.POST("/stress/mem/stop", stressHandler.MemStop)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
