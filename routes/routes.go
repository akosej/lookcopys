package routes

import (
	"github.com/akosej/lookcopys/controllers"
	"github.com/gofiber/fiber/v2"
)

func RoutesApi(app *fiber.App) {
	//  ---- Authentication routes
	api := app.Group("/api")
	api.Get("/status", controllers.ApiStatus).Name("Status")
}
