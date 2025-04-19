package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ReqPayMentTruemoney struct {
	Link string `json:"link"`
}

type RedeemRequest struct {
	Mobile      string `json:"mobile"`
	VoucherHash string `json:"voucher_hash"`
}

func PayMentTruemoney(c *fiber.Ctx) error {
	var req ReqPayMentTruemoney

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถอ่านข้อมูลได้: " + err.Error(),
		})
	}

	fmt.Println(req)

	isValid := strings.HasPrefix(req.Link, "https://gift.truemoney.com/")
	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ลิงก์ไม่ถูกต้อง กรุณาใส่ลิงก์ซองอั่งเปา",
		})
	}

	parsedUrl, err := url.Parse(req.Link)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ลิงก์ไม่ถูกต้อง",
		})
	}

	voucher := parsedUrl.Query().Get("v")

	redeemUrl := "https://gift.truemoney.com/campaign/vouchers/" + voucher + "/redeem"

	body := RedeemRequest{
		Mobile:      "0650269486",
		VoucherHash: voucher,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "แปลง body ไม่สำเร็จ: " + err.Error(),
		})
	}

	reqHttp, err := http.NewRequest("POST", redeemUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "สร้างคำขอไม่สำเร็จ: " + err.Error(),
		})
	}
	reqHttp.Header.Set("Content-Type", "application/json")
	reqHttp.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(reqHttp)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Redeem ล้มเหลว: " + err.Error(),
		})
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ไม่สามารถแปลงข้อมูล: " + err.Error(),
		})
	}

	status, ok := result["status"].(map[string]interface{})
	if !ok {
		return c.JSON(fiber.Map{
			"success": false,
		})
	}

	code, ok := status["code"].(string)
	if !ok {
		return c.JSON(fiber.Map{
			"success": false,
		})
	}

	if code == "VOUCHER_EXPIRED" || code == "TARGET_USER_REDEEMED" || code == "VOUCHER_OUT_OF_STOCK" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "ลิ้งนี้ถูกใช้แล้ว หรือ หมดอายุ",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})

}
