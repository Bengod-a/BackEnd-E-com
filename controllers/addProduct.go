package controllers

import (
	"fmt"
	"strconv"

	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/gofiber/fiber/v2"
)

type ReqProduct struct {
	Name       string   `form:"name"`
	Price      string   `form:"price"`
	CategoryID string   `form:"categoryId"`
	Quantity   string   `form:"quantity"`
	UserID     string   `form:"userId"`
	Images     []string `json:"images"`
}

func AddProduct(c *fiber.Ctx) error {
	var req ReqProduct
	fmt.Println(req)

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid file upload",
		})
	}

	req.Name = form.Value["name"][0]
	req.Price = form.Value["price"][0]
	req.CategoryID = form.Value["categoryId"][0]
	req.UserID = form.Value["userId"][0]
	req.Quantity = form.Value["quantity"][0]

	if req.Name == "" || req.Price == "" || req.CategoryID == "" || req.UserID == "" || req.Quantity == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "กรุณากรอกข้อมูลให้ครบ",
		})
	}

	fmt.Println(form)

	priceFloat, err := strconv.ParseFloat(req.Price, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ราคาสินค้าต้องเป็นตัวเลข",
		})
	}

	categoryID, err := strconv.ParseUint(req.CategoryID, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID หมวดหมู่ต้องเป็นตัวเลข",
		})
	}

	var category models.Category
	if err := db.DB.First(&category, categoryID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "ไม่พบหมวดหมู่ที่เลือก",
		})
	}

	var images []models.Images
	for _, imageURL := range form.Value["images"] {
		image := models.Images{
			URL: imageURL,
		}
		images = append(images, image)
	}

	qty, err := strconv.Atoi(req.Quantity)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "quantity format",
		})
	}

	product := models.Product{
		Name:       req.Name,
		Price:      priceFloat,
		Quantity:   qty,
		Categories: []models.Category{category},
		Images:     images,
	}

	db.DB.Create(&product)

	return c.JSON(fiber.Map{
		"message": "Product added successfully",
		"data":    product,
	})
}
