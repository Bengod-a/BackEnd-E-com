package controllers

import (
	"errors"

	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Reqcategory struct {
	CategoryName string `json:"categoryName"`
	Icon         string `json:"icon"`
}

func Addcategory1(c *fiber.Ctx) error {
	var req Reqcategory

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if req.CategoryName == "" || req.Icon == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "กรุณากรอกข้อมูลให้ครบ",
		})
	}

	var existingCategory models.Category1
	err := db.DB.Where("name = ?", req.CategoryName).First(&existingCategory).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Database error: " + err.Error(),
		})
	} else if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "มีชื่อนี้แล้ว",
		})
	}

	category := models.Category1{Name: req.CategoryName, Icon: req.Icon}
	if err := db.DB.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in Create" + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "สมัครสมาชิกสำเร็จ",
		"data":    category,
	})

}

func Addcategory(c *fiber.Ctx) error {
	var req Reqcategory

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if req.CategoryName == "" || req.Icon == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "กรุณากรอกข้อมูลให้ครบ",
		})
	}

	var existingCategory models.Category2
	err := db.DB.Where("name = ?", req.CategoryName).First(&existingCategory).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Database error: " + err.Error(),
		})
	} else if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "มีชื่อนี้แล้ว",
		})
	}

	category := models.Category2{Name: req.CategoryName, Icon: req.Icon}
	if err := db.DB.Create(&category).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error in Create" + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "สมัครสมาชิกสำเร็จ",
		"data":    category,
	})

}
