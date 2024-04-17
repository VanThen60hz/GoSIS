package routes

import (
	controllers "GoSIS/controller"

	"github.com/gofiber/fiber/v2"
)

func GetRoute(app *fiber.App) {
	// All routes related to users comes here

	// Employee route
	app.Get("/employee", controllers.GetAllEmployees)

	// Personal route
	app.Get("/personal", controllers.GetAllPersonals)
	app.Post("/personal", controllers.CreatePersonal)

	// Merge Employee-Personal
	app.Get("/merge-person", controllers.MergeData)

	// Dashboard
	app.Get("/gender-ratio", controllers.GenderRatio)
}
