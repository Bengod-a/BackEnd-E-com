package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
)

type PaymentRequest struct {
	Amount int64 `json:"amount"`
}

func CreatePaymentIntent(c *fiber.Ctx) error {
	var req PaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(req.Amount),
		Currency: stripe.String(string(stripe.CurrencyTHB)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Payment failed", "details": err.Error()})
	}

	return c.JSON(fiber.Map{
		"clientSecret": pi.ClientSecret,
	})
}
