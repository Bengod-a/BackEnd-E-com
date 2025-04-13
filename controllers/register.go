package controllers

import (
	"strings"

	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type Req struct {
	Username    string `json:"username"`
	Id          int    `json:"id"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Phonenumber string `json:"phonenumber"`
	Conpassword string `json:"conpassword"`
}

func Register(c *fiber.Ctx) error {
	var req Req

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถอ่านข้อมูลได้",
		})
	}

	if strings.TrimSpace(req.Username) == "" ||
		strings.TrimSpace(req.Lastname) == "" ||
		strings.TrimSpace(req.Email) == "" ||
		strings.TrimSpace(req.Phonenumber) == "" ||
		strings.TrimSpace(req.Password) == "" ||
		strings.TrimSpace(req.Conpassword) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "กรุณากรอกข้อมูลให้ครบ",
		})
	}

	if req.Password != req.Conpassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "รหัสผ่านและยืนยันรหัสผ่านไม่ตรงกัน",
		})
	}

	var existingUser models.User
	err := db.DB.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "อีเมลนี้ถูกใช้งานแล้ว",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถเข้ารหัสรหัสผ่านได้",
		})
	}

	user := models.User{
		Name:        req.Username,
		Lastname:    &req.Lastname,
		Email:       req.Email,
		Password:    string(hashedPassword),
		Phonenumber: &req.Phonenumber,
		Role:        models.UserRole,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถสมัครสมาชิกได้",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "สมัครสมาชิกสำเร็จ",
		"data":    user,
	})
}

func Edituser(c *fiber.Ctx) error {
	var req Req
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถอ่านข้อมูลได้",
		})
	}

	if strings.TrimSpace(req.Username) == "" ||
		strings.TrimSpace(req.Lastname) == "" ||
		strings.TrimSpace(req.Email) == "" ||
		strings.TrimSpace(req.Phonenumber) == "" {
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
			"message": "ไม่มี ID นี้",
		})
	}

	user.Name = req.Username
	user.Lastname = &req.Lastname
	user.Email = req.Email
	user.Phonenumber = &req.Phonenumber

	if err := db.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "เกิดข้อผิดพลาดในการบันทึกข้อมูล",
		})
	}

	return c.JSON(fiber.Map{
		"message": "อัปเดตข้อมูลส่วนตัวสำเร็จ",
		"data":    user,
	})
}
