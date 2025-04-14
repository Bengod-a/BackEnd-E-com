package controllers

import (
	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/gofiber/fiber/v2"
)

func GetProductsid(c *fiber.Ctx) error {
	id := c.Params("id")

	var product []models.Product

	if err := db.DB.Preload("Categories1").Preload("Categories2").Preload("Images").First(&product, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Error in GetProductsid",
		})
	}

	return c.JSON(fiber.Map{
		"data": product,
	})

}
