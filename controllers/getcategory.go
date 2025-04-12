package controllers

import (
	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/gofiber/fiber/v2"
)

func Getcategory(c *fiber.Ctx) error {

	var category []models.Category

	err := db.DB.Find(&category).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in Getcategory",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    category,
	})
}
