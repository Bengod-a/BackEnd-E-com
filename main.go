package main

import (
	"github.com/Bengod-a/DB-GO/controllers"
	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	db.InitDB()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	app.Post("/products", middleware.AdminCheck, controllers.AddProduct)
	app.Get("/products", middleware.AdminCheck, controllers.GetProduct)
	app.Get("/products/:id", middleware.AdminCheck, controllers.GetProductsid)
	app.Post("/category", middleware.AdminCheck, controllers.Addcategory)
	app.Post("/uploadimage", middleware.AdminCheck, controllers.UploadImages)
	app.Post("/removeimage", middleware.AdminCheck, controllers.HandleRemoveImage)
	app.Get("/category", middleware.AdminCheck, controllers.Getcategory)
	app.Get("/user", middleware.AdminCheck, controllers.Getuser)

	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)

	if err := app.Listen(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}
