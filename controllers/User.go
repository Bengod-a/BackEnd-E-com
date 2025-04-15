package controllers

import (
	"strings"

	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/gofiber/fiber/v2"
)

type ReqCreateAddress struct {
	Id       int    `json:"id"`
	Address  string `json:"address"`
	Province string `json:"province"`
	Amphure  string `json:"amphure"`
	Tambon   string `json:"tambon"`
	Zipcode  int    `json:"zipcode"`
}

func CreateAddress(c *fiber.Ctx) error {
	var req ReqCreateAddress
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถอ่านข้อมูลได้",
		})
	}

	if strings.TrimSpace(req.Province) == "" ||
		strings.TrimSpace(req.Amphure) == "" ||
		strings.TrimSpace(req.Address) == "" ||
		strings.TrimSpace(req.Tambon) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "กรุณากรอกข้อมูลให้ครบ",
		})
	}

	var user models.User
	err := db.DB.Where("id = ?", req.Id).First(&user).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Error",
		})
	}

	address := models.Address{
		Address:  req.Address,
		Province: req.Province,
		Amphure:  req.Amphure,
		Tambon:   req.Tambon,
		Zipcode:  req.Zipcode,
		UserID:   uint(req.Id),
	}

	if err := db.DB.Create(&address).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Create Address Errro",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "เพิ่มที่อยู่สำเร็จ",
		"data":    address,
	})

}

func DeleteAddress(c *fiber.Ctx) error {
	id := c.Params("id")

	var address models.Address
	if err := db.DB.Where("id = ?", id).First(&address).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "ไม่พบที่อยู่",
		})
	}

	if err := db.DB.Unscoped().Delete(&address).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Delete Address Error",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ลบที่อยู่สำเร็จ",
		"data":    address,
	})

}
