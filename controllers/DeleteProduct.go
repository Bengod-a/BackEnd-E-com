package controllers

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
)

type ReqDeleteProduct struct {
	Urls []string `json:"urls"`
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	productID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var req ReqDeleteProduct

	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		log.Println("Cloudinary config error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Cloudinary configuration failed",
		})
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var imageURL []string = req.Urls

	for _, url := range imageURL {
		parts := strings.Split(url, "/")
		publicID := strings.Split(parts[len(parts)-1], ".")[0]

		_, err = cld.Upload.Destroy(context.Background(), uploader.DestroyParams{
			PublicID: publicID,
		})

	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ลบรูปไม่สำเร็จ",
		})
	}

	var product models.Product

	if err := db.DB.First(&product, "id = ?", productID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	if err := db.DB.Unscoped().Delete(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error delete product",
		})
	}

	return c.JSON(fiber.Map{
		"message": "ลบสินค้าสำเร็จ",
	})
}
