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

type ReqAddToCart struct {
	Productid int `json:"Productid"`
	Price     int `json:"price"`
	Quantity  int `json:"quantity"`
	Count     int `json:"count"`
	Id        int `json:"id"`
}

func AddToCart(c *fiber.Ctx) error {
	var req ReqAddToCart
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถอ่านข้อมูลได้",
		})
	}

	if int(req.Productid) == 0 || int(req.Quantity) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ข้อมูลไม่ครบ",
		})
	}
	var product models.Product
	if err := db.DB.Where("id = ?", req.Productid).First(&product).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่พบสินค้า",
		})
	}

	allcount := float64(req.Quantity) + float64(req.Count)

	if int(allcount) > product.Quantity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "สินค้าไม่พอ",
		})
	}

	var cart models.Cart
	err := db.DB.Where("order_by_id = ?", req.Id).First(&cart).Error
	if err != nil {
		cart = models.Cart{
			CartTotal: float64(req.Price * req.Quantity),
			OrderByID: uint(req.Id),
		}

		if err := db.DB.Create(&cart).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "ไม่สามารถสร้างตะกร้าได้",
			})
		}

	} else {
		cart.CartTotal += float64(req.Price * req.Quantity)
		if err := db.DB.Save(&cart).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "ไม่สามารถอัปเดตตะกร้าได้",
			})
		}
	}

	var productOnCart models.ProductOnCart

	err = db.DB.Where("product_id = ? AND cart_id = ?", req.Productid, cart.ID).First(&productOnCart).Error
	if err != nil {
		productOnCart = models.ProductOnCart{
			CartID:    cart.ID,
			ProductID: uint(req.Productid),
			Count:     req.Quantity,
			Price:     float64(req.Price),
		}
		if err := db.DB.Create(&productOnCart).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "เพิ่มสินค้าในตะกร้าไม่สำเร็จ",
			})
		}
	} else {
		productOnCart.Count += req.Quantity
		productOnCart.Price = float64(req.Price)
		if err := db.DB.Save(&productOnCart).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "อัปเดตสินค้าในตะกร้าไม่สำเร็จ",
			})
		}
	}
	db.DB.
		Preload("Cart").
		Preload("Product").
		Preload("Product.Images").
		Preload("Product.Categories1").
		Preload("Product.Categories2").
		First(&productOnCart, productOnCart.ID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "เพิ่มสินค้าในตะกร้าสำเร็จ",
		"cart":    cart,
		"product": productOnCart,
	})
}

func DeleteproductOnCart(c *fiber.Ctx) error {
	id := c.Params("id")

	var productoncart models.ProductOnCart
	if err := db.DB.Where("id = ?", id).First(&productoncart).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่พบตระกร้า",
		})
	}

	if err := db.DB.Unscoped().Delete(&productoncart).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Delete ProductOnCart Error",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "ลบสินค้าจากตะกร้าเรียบร้อยแล้ว",
	})
}
