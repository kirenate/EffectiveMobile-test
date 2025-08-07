package presentations

import (
	"github.com/gofiber/fiber/v2"
	"main.go/services"
)

func (r *Presentation) postSubscription(c *fiber.Ctx) error {
	var req *services.SubscriptionRequest

	err := c.BodyParser(&req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}
	err = r.service.ProcessSubscriptionRequest(c.UserContext(), req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: err.Error()}
	}
	return c.JSON(fiber.Map{"status": "success"})
}
