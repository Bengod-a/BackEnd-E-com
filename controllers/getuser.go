package controllers

import (
	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/gofiber/fiber/v2"
)

func Getuser(c *fiber.Ctx) error {

	var user models.User
	err := db.DB.Find(&user).Error

	if err != nil {
		return c.JSON(fiber.StatusBadRequest)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"name":    user.Name,
			"email":   user.Email,
			"role":    user.Role,
			"order":   user.Orders,
			"address": user.Addresses,
		},
	})
}
