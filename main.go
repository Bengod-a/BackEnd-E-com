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
		AllowMethods:     "GET,POST,PUT,DELETE, PATCH",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	// DropAllTables()

	// Admin
	app.Post("/products", middleware.AdminCheck, controllers.AddProduct)
	app.Get("/products", middleware.AdminCheck, controllers.GetProduct)
	app.Get("/products/:id", middleware.AdminCheck, controllers.GetProductsid)
	app.Patch("/products/:id", controllers.ChangProduct)
	app.Delete("/products/:id", middleware.AdminCheck, controllers.DeleteProduct)
	app.Post("/category", middleware.AdminCheck, controllers.Addcategory)
	app.Get("/category", middleware.AdminCheck, controllers.Getcategory)
	app.Get("/category/:id", middleware.AdminCheck, controllers.GetcategoryId)
	app.Patch("/category/:id", middleware.AdminCheck, controllers.Editcategory)
	app.Post("/uploadimage", middleware.AdminCheck, controllers.UploadImages)
	app.Post("/removeimage", middleware.AdminCheck, controllers.HandleRemoveImage)
	app.Post("/removeimageinproduct/:id", middleware.AdminCheck, controllers.RemoveImageInProduct)
	app.Get("/user", middleware.AdminCheck, controllers.Getuser)
	// Admin

	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)

	if err := app.Listen(":8080"); err != nil {
		panic("Failed to start server: " + err.Error())
	}
}

func DropAllTables() {
	db.DB.Exec(`DO $$ DECLARE
		r RECORD;
	BEGIN
		FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
			EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
		END LOOP;
	END $$;`)
}
