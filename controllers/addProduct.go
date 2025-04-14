// package controllers

// import (
// 	"encoding/json"
// 	"fmt"
// 	"strconv"

// 	"github.com/Bengod-a/DB-GO/db"
// 	"github.com/Bengod-a/DB-GO/models"
// 	"github.com/gofiber/fiber/v2"
// 	"gorm.io/datatypes"
// )

// type ReqProduct struct {
// 	Name       string   `form:"name"`
// 	Price      string   `form:"price"`
// 	CategoryID string   `form:"categoryId"`
// 	Quantity   string   `form:"quantity"`
// 	UserID     string   `form:"userId"`
// 	Images     []string `json:"images"`
// 	Key        string   `json:"key"`
// 	Value      string   `json:"value"`
// }

// func AddProduct(c *fiber.Ctx) error {
// 	var req ReqProduct

// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Invalid file upload",
// 		})
// 	}

// 	req.Name = form.Value["name"][0]
// 	req.Price = form.Value["price"][0]
// 	req.CategoryID = form.Value["categoryId"][0]
// 	req.UserID = form.Value["userId"][0]
// 	req.Quantity = form.Value["quantity"][0]
// 	req.Key = form.Value["key"][0]
// 	req.Value = form.Value["value"][0]
// 	fmt.Println(req)

// 	if req.Name == "" || req.Price == "" || req.CategoryID == "" || req.UserID == "" || req.Quantity == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "กรุณากรอกข้อมูลให้ครบ",
// 		})
// 	}

// 	priceFloat, err := strconv.ParseFloat(req.Price, 64)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "ราคาสินค้าต้องเป็นตัวเลข",
// 		})
// 	}

// 	categoryID, err := strconv.ParseUint(req.CategoryID, 10, 32)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "ID หมวดหมู่ต้องเป็นตัวเลข",
// 		})
// 	}

// 	var category models.Category
// 	if err := db.DB.First(&category, categoryID).Error; err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"message": "ไม่พบหมวดหมู่ที่เลือก",
// 		})
// 	}

// 	qty, err := strconv.Atoi(req.Quantity)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"success": false,
// 			"message": "quantity format",
// 		})
// 	}

// 	var images []models.Images
// 	for _, imageURL := range form.Value["images"] {
// 		image := models.Images{
// 			URL: imageURL,
// 		}
// 		images = append(images, image)
// 	}

// 	specsData := map[string]string{
// 		req.Key: req.Value,
// 	}

// 	specsJSON, err := json.Marshal(specsData)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "เกิดข้อผิดพลาดในการแปลงข้อมูล JSON",
// 		})
// 	}

// 	product := models.Product{
// 		Name:       req.Name,
// 		Price:      priceFloat,
// 		Quantity:   qty,
// 		Categories: []models.Category{category},
// 		Images:     images,
// 		Specs:      datatypes.JSON(specsJSON),
// 	}

// 	db.DB.Create(&product)

// 	return c.JSON(fiber.Map{
// 		"message": "เพิ่มสินค้าสำเร็จ",
// 		"data":    product,
// 	})
// }

package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/Bengod-a/DB-GO/db"
	"github.com/Bengod-a/DB-GO/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

type ReqProduct struct {
	Name        string   `form:"name"`
	Price       string   `form:"price"`
	CategoryID1 string   `form:"categoryId1"`
	CategoryID2 string   `form:"categoryId2"`
	Quantity    string   `form:"quantity"`
	UserID      string   `form:"userId"`
	Images      []string `json:"images"`
	Key         string   `json:"key"`
	Value       string   `json:"value"`
}

func AddProduct(c *fiber.Ctx) error {
	var req ReqProduct

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid file upload",
		})
	}

	req.Name = getFirstValue(form.Value["name"])
	req.Price = getFirstValue(form.Value["price"])
	req.CategoryID1 = getFirstValue(form.Value["categoryId1"])
	req.CategoryID2 = getFirstValue(form.Value["categoryId2"])
	req.UserID = getFirstValue(form.Value["userId"])
	req.Quantity = getFirstValue(form.Value["quantity"])

	if req.Name == "" || req.Price == "" || req.CategoryID1 == "" || req.CategoryID2 == "" || req.UserID == "" || req.Quantity == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "กรุณากรอกข้อมูลให้ครบ",
		})
	}

	priceFloat, err := strconv.ParseFloat(req.Price, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ราคาสินค้าต้องเป็นตัวเลข",
		})
	}

	categoryID1, err := strconv.ParseUint(req.CategoryID1, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID หมวดหมู่ต้องเป็นตัวเลข",
		})
	}

	categoryID2, err := strconv.ParseUint(req.CategoryID2, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "ID หมวดหมู่ต้องเป็นตัวเลข",
		})
	}

	var category1 models.Category1
	if err := db.DB.First(&category1, categoryID1).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "ไม่พบหมวดหมู่ที่เลือก",
		})
	}

	var category2 models.Category2
	if err := db.DB.First(&category2, categoryID2).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "ไม่พบหมวดหมู่ที่เลือก",
		})
	}

	qty, err := strconv.Atoi(req.Quantity)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "quantity format",
		})
	}

	var images []models.Images
	for _, imageURL := range form.Value["images"] {
		image := models.Images{
			URL: imageURL,
		}
		images = append(images, image)
	}

	specsData := form.Value["envs"]
	var specs []map[string]string
	if err := json.Unmarshal([]byte(specsData[0]), &specs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "เกิดข้อผิดพลาดในการแปลงข้อมูล",
		})
	}

	specsJSON, err := json.Marshal(specs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "เกิดข้อผิดพลาดในการแปลงข้อมูล JSON",
		})
	}

	product := models.Product{
		Name:        req.Name,
		Price:       priceFloat,
		Quantity:    qty,
		Categories1: []models.Category1{category1},
		Categories2: []models.Category2{category2},
		Images:      images,
		Specs:       datatypes.JSON(specsJSON),
	}

	db.DB.Create(&product)

	return c.JSON(fiber.Map{
		"message": "เพิ่มสินค้าสำเร็จ",
		"data":    product,
	})
}

func getFirstValue(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}
