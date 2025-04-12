package controllers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func HandleRemoveImage(c *fiber.Ctx) error {
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

	var req struct {
		Urls []string `json:"urls"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	fmt.Println(req)

	if len(req.Urls) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No image URL",
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

	return c.JSON(fiber.Map{
		"message": "ลบรูปสำเร็จ",
	})
}
