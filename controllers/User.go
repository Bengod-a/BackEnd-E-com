package controllers

import (
	"fmt"
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

type ReqAddToFavorite struct {
	Productid int `json:"productid"`
	Id        int `json:"id"`
}

func AddToFavorite(c *fiber.Ctx) error {
	var req ReqAddToFavorite

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถอ่านข้อมูลได้",
		})
	}

	if int(req.Productid) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ข้อมูลไม่ครบ",
		})
	}

	var product models.Product
	if err := db.DB.Where("id = ?", req.Productid).First(&product).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่พบ Product",
		})
	}

	favorite := models.Favorite{
		UserID:    uint(req.Id),
		ProductID: uint(req.Productid),
	}
	if err := db.DB.Create(&favorite).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Error Create Favorite",
		})
	}

	if err := db.DB.
		Where("user_id = ? AND product_id = ?", req.Id, req.Productid).
		First(&favorite).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถดึงข้อมูล favorite พร้อม product ได้",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "เพิ่มในรายการโปรดแล้ว",
		"data":    favorite,
	})

}

func DeleteFavorite(c *fiber.Ctx) error {
	var req ReqAddToFavorite

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถอ่านข้อมูลได้",
		})
	}

	if int(req.Productid) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ข้อมูลไม่ครบ",
		})
	}

	var product models.Product
	if err := db.DB.Where("id = ?", req.Productid).First(&product).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่พบ Product",
		})
	}

	var favorite models.Favorite
	if err := db.DB.Where("product_id = ? AND user_id = ?", req.Productid, req.Id).First(&favorite).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่พบ รายการโปรด",
		})
	}

	if err := db.DB.Unscoped().Delete(&favorite).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถลบรายการโปรดได้",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "!ลบรายการโปรดแล้ว",
		"data":    favorite,
	})

}

func GetFavorite(c *fiber.Ctx) error {
	favoriteID := c.Query("favoriteid")
	userID := c.Query("id")

	favoriteIDStrings := strings.Split(favoriteID, ",")

	var favorite []models.Favorite
	if err := db.DB.
		Preload("Product").
		Preload("Product.Images").
		Where("id IN (?) AND user_id = ?", favoriteIDStrings, userID).
		Find(&favorite).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    favorite,
	})

}

type StripePaymentIntent struct {
	StripeId        string `json:"stripeId"`
	LinkTruemoney   string `json:"linkTruemoney"`
	Currency        string `json:"currency"`
	OptionsPayMent  string `json:"optionspayment"`
	Id              int    `json:"id"`
	Count           []int  `json:"count"`
	ProductOnCartID []int  `json:"productoncartid"`
	Productid       []int  `json:"productid"`
	Addressid       int    `json:"addressid"`
	Amount          int64  `json:"amount"`
	Status          string `json:"status"`
}

func SaveOrders(c *fiber.Ctx) error {
	var req StripePaymentIntent

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถอ่านข้อมูลได้",
		})
	}

	fmt.Println(req.Count)

	var user models.User
	if err := db.DB.Where("id = ?", req.Id).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่พบ User",
		})
	}

	var products []models.Product
	if err := db.DB.Where("id IN ?", req.Productid).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
		})
	}

	addressID := uint(req.Addressid)
	Amount := int(req.Amount) / 100
	order := models.Order{
		PayMentID:      req.StripeId,
		LinkTruemoney:  req.LinkTruemoney,
		OptionsPayMent: req.OptionsPayMent,
		Currency:       &req.Currency,
		PayMentStatus:  req.Status,
		AddressID:      &addressID,
		Status:         "กำลังดำเนินการ",
		Amount:         int(Amount),
		OrderByID:      uint(req.Id),
	}
	if err := db.DB.Create(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Create Orders Error",
		})
	}

	count := req.Count
	productid := req.Productid
	productoncartid := req.ProductOnCartID
	if len(count) != len(productid) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
		})
	}

	for i := 0; i < len(count); i++ {
		productId := productid[i]
		productCount := count[i]
		productoncartid := productoncartid[i]

		for _, item := range products {
			if int(item.ID) == productId {
				if item.Quantity < productCount {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"success": false,
						"message": "สินค้าไม่พอ",
					})
				}

				item.Quantity -= productCount
				if err := db.DB.Save(&item).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"success": false,
						"message": "ลด stock ไม่ได้",
					})
				}

				productOnOrder := models.ProductOnOrder{
					ProductID: uint(productid[i]),
					OrderID:   order.ID,
					Count:     productCount,
				}
				if err := db.DB.Create(&productOnOrder).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"success": false,
						"message": "Error in CreateproductOnOrder" + err.Error(),
					})
				}

				var productoncart models.ProductOnCart
				if err := db.DB.First(&productoncart, "id = ?", productoncartid).Error; err != nil {
					return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"success": false,
						"error":   "Error Delete ProductOnCart",
					})
				}

				if err := db.DB.Unscoped().Delete(&productoncart).Error; err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"success": false,
						"message": "ลบ productOnCart ไม่ได้",
					})
				}

				break
			}
		}
	}

	if err := db.DB.
		Preload("OrderBy").
		Preload("Address").
		First(&order, order.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    order,
	})

}
