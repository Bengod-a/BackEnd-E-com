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
