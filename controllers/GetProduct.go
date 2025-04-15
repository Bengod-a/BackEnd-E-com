package controllers

import (
	"strconv"

	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/gofiber/fiber/v2"
)

func GetProduct(c *fiber.Ctx) error {
	var products []models.Product
	var total int64

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("limit", "1"))
	offset := (page - 1) * pageSize

	db.DB.Model(&models.Product{}).Count(&total)

	err := db.DB.Preload("Categories2").Preload("Categories1").Preload("Images").Limit(pageSize).Offset(offset).Find(&products).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in GetProduct",
		})
	}

	return c.JSON(fiber.Map{
		"data": products,
	})

}

func GetBesSaleProduct(c *fiber.Ctx) error {
	var products []models.Product

	if err := db.DB.Preload("Categories1").Preload("Categories2").Preload("Images").Order("sold DESC").Limit(10).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Error in GetBesSaleProduct",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ดึงสินค้าขายดีสำเร็จ",
		"data":    products,
	})
}

func GetProductById(c *fiber.Ctx) error {
	id := c.Params("id")
	var products []models.Product

	err := db.DB.Preload("Categories2").Preload("Categories1").Preload("Images").Where("id = ?", &id).Find(&products).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in GetProductById",
		})
	}

	return c.JSON(fiber.Map{
		"data": products,
	})

}
