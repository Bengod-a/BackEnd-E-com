package main

import (
	"log"
	"os"

	"github.com/Bengod-a/DB-GO/controllers"
	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v82"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	stripe.Key = stripeKey

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
	// category
	// หมวดหมู่หลัก
	app.Post("/category1", middleware.AdminCheck, controllers.Addcategory1)
	app.Get("/category1", middleware.AdminCheck, controllers.Getcategory1)
	// หมวดหมู่รอง
	app.Post("/category", middleware.AdminCheck, controllers.Addcategory)
	app.Get("/category", middleware.AdminCheck, controllers.Getcategory)

	app.Get("/category/:id", middleware.AdminCheck, controllers.GetcategoryId)
	app.Patch("/category/:id", middleware.AdminCheck, controllers.Editcategory)
	// Image
	app.Post("/uploadimage", middleware.AdminCheck, controllers.UploadImages)
	app.Post("/removeimage", middleware.AdminCheck, controllers.HandleRemoveImage)
	app.Post("/removeimageinproduct/:id", middleware.AdminCheck, controllers.RemoveImageInProduct)
	// GetUser
	app.Get("/user", middleware.AdminCheck, controllers.Getuser)
	// Admin

	// PayMent
	app.Post("/create-payment-intent", middleware.UserCheck, controllers.CreatePaymentIntent)
	app.Post("/saveorders", middleware.UserCheck, controllers.SaveOrders)
	app.Post("/paymenttruemoney", controllers.PayMentTruemoney)

	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)
	app.Post("/address", middleware.UserCheck, controllers.CreateAddress)
	app.Delete("/address/:id", middleware.UserCheck, controllers.DeleteAddress)
	app.Patch("/edituser", middleware.UserCheck, controllers.Edituser)
	app.Delete("/productoncart/:id", middleware.UserCheck, controllers.DeleteproductOnCart)
	app.Post("/addtocart", middleware.UserCheck, controllers.AddToCart)

	app.Post("/addtofavorite", middleware.UserCheck, controllers.AddToFavorite)
	app.Get("/favorite", controllers.GetFavorite)
	app.Delete("/deletefavorite", middleware.UserCheck, controllers.DeleteFavorite)

	app.Get("/GetBesSaleProduct", controllers.GetBesSaleProduct)
	app.Get("/product/:id", controllers.GetProductById)

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
