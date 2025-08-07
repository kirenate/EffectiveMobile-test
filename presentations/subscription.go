package presentations

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
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

func (r *Presentation) getSubscription(c *fiber.Ctx) error {
	var req *services.SubscriptionRequest

	err := c.BodyParser(&req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	sub, err := r.service.ProcessSubscriptionGetRequest(c.UserContext(), req)
	if err != nil {
		return errors.Wrap(err, "failed to process subscription get request")
	}

	return c.JSON(sub)
}

func (r *Presentation) deleteSubscription(c *fiber.Ctx) error {
	var req *services.SubscriptionRequest

	err := c.BodyParser(&req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = r.service.ProcessSubscriptionDeleteRequest(c.UserContext(), req)
	if err != nil {
		return errors.Wrap(err, "failed to process subscription delete request")
	}

	return c.JSON(fiber.Map{"status": "success"})
}
