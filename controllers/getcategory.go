package controllers

import (
	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/gofiber/fiber/v2"
)

func Getcategory1(c *fiber.Ctx) error {

	var category []models.Category1

	err := db.DB.Preload("Products").Find(&category).Error
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

func Getcategory(c *fiber.Ctx) error {

	var category []models.Category2

	err := db.DB.Preload("Products").Find(&category).Error
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

func GetcategoryId(c *fiber.Ctx) error {
	id := c.Params("id")

	var category models.Category2
	err := db.DB.Preload("Products.Images").First(&category, "id = ?", id).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in GetCategoryId",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    category,
	})
}

func Editcategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var req Reqcategory

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in BodyParser",
		})
	}
	var category models.Category2
	if err := db.DB.Where("id = ?", id).First(&category).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in Editcategory",
		})
	}

	category.Name = req.CategoryName
	category.Icon = req.Icon

	if err := db.DB.Save(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error update category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": category,
	})

}
