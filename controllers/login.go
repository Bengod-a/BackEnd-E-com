package controllers

import (
	"time"

	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/Bengod-a/DB-GO/services"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ReqLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var req ReqLogin

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถอ่านข้อมูลได้: " + err.Error(),
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "กรุณากรอกข้อมูลให้ครบ",
		})
	}

	var user models.User
	if err := db.DB.Preload("Addresses").Preload("Carts.Products").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "ไม่พบมีเมลนี้",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการค้นหาผู้ใช้",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "อีเมลหรือรหัสผ่านไม่ถูกต้อง",
		})
	}

	token, err := services.JWTAuthService().GenerateToken(int(user.ID), user.Email, string(user.Role), 6)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "JWT Gen Token Error",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:    "Benjwt",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 6),
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "ล็อกอินสำเร็จ",
		"token":   token,
		"user": fiber.Map{
			"id":          user.ID,
			"username":    user.Name,
			"email":       user.Email,
			"role":        user.Role,
			"lastname":    user.Lastname,
			"phonenumber": user.Phonenumber,
			"orders":      user.Orders,
			"address":     user.Addresses,
			"favorite":    user.Favorites,
			"carts":       user.Carts,
		},
	})

}
