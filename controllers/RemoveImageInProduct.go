package controllers

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var req struct {
	Urls []string `json:"urls"`
}

func RemoveImageInProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	err := godotenv.Load()
	if err != nil {
		log.Println("Error .env")

	}

	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "connect Cloudinary Error",
		})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	imageURL := req.Urls[0]
	parts := strings.Split(imageURL, "/")
	publicID := strings.Split(parts[len(parts)-1], ".")[0]

	_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ลบรูปไม่สำเร็จ",
		})
	}

	var image models.Images

	if err := db.DB.First(&image, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "image not found",
		})
	}

	if err := db.DB.Unscoped().Delete(&image).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error delete product",
		})
	}

	return c.JSON(fiber.Map{
		"message": "ลบรูปสำเร็จ",
	})

}
