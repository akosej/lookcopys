package main

import (
	"github.com/akosej/lookcopys/routes"
	"github.com/akosej/lookcopys/system"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"time"
)

func init() {
	system.CreateDirectoryIfDoesntExist(system.Path + "/records")
	go system.MonitorUsb()
	system.SendNotifyDesktop("Alerta", "LookCopys se ha iniciado")
}
func main() {
	// ------------------------------------------------------------------------------------------------------------
	// ------- RUN API-AGA-SENTINEL
	// ------------------------------------------------------------------------------------------------------------
	app := fiber.New(fiber.Config{
		AppName:       "LookCopys v1.0",
		CaseSensitive: true,
		//EnablePrintRoutes: true,
		//GETOnly: true,
	})

	app.Use(cors.New(cors.Config{
		AllowCredentials: false,
	}))

	app.Static("/", system.Path+"/frontend", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	})

	routes.RoutesApi(app)
	_ = app.Listen(":5323")

}
