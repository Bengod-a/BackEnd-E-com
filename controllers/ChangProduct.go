package controllers

import (
	"fmt"
	"strconv"

	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/gofiber/fiber/v2"
)

func ChangProduct(c *fiber.Ctx) error {
	var req ReqProduct
	fmt.Println(req)

	id := c.Params("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error ProductID",
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ฟอร์มไม่ถูกต้อง",
		})
	}

	req.Name = form.Value["name"][0]
	req.Price = form.Value["price"][0]
	req.CategoryID = form.Value["categoryId"][0]
	req.Quantity = form.Value["quantity"][0]

	priceFloat, _ := strconv.ParseFloat(req.Price, 64)
	categoryID, _ := strconv.ParseUint(req.CategoryID, 10, 32)
	qty, _ := strconv.Atoi(req.Quantity)

	var product models.Product
	if err := db.DB.Preload("Images").Preload("Categories").First(&product, productID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "ไม่พบสินค้า",
		})
	}

	db.DB.Where("product_id = ?", product.ID).Delete(&models.Images{})

	var images []models.Images
	for _, url := range form.Value["images"] {
		images = append(images, models.Images{URL: url})
	}

	var category models.Category
	if err := db.DB.First(&category, categoryID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "ไม่พบหมวดหมู่",
		})
	}

	product.Name = req.Name
	product.Price = priceFloat
	product.Quantity = qty
	product.Images = images
	product.Categories = []models.Category{category}

	db.DB.Save(&product)

	return c.JSON(fiber.Map{
		"message": "อัปเดตสินค้าสำเร็จ",
		"data":    product,
	})
}
