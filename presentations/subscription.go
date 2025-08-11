package presentations

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"main.go/services"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())

func (r *Presentation) postSubscription(c *fiber.Ctx) error {
	var req *services.SubscriptionRequest

	err := c.BodyParser(&req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = Validate.Struct(req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = r.service.ProcessSubscriptionRequest(c.UserContext(), req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: err.Error()}
	}

	zerolog.Ctx(c.UserContext()).
		Info().
		Interface("subscription_request", req).
		Msg("new.subscription.request")

	return c.JSON(fiber.Map{"status": "success"})
}

func (r *Presentation) getSubscription(c *fiber.Ctx) error {
	sub, err := r.service.ProcessSubscriptionGetRequest(c.UserContext())
	if err != nil {
		return errors.Wrap(err, "failed to process subscription get request")
	}

	zerolog.Ctx(c.UserContext()).
		Info().
		Str("subscription_request", string(c.Body())).
		Msg("new.get.subscription.request")

	return c.JSON(sub)
}

func (r *Presentation) deleteSubscription(c *fiber.Ctx) error {
	var req *services.SubscriptionRequest

	err := c.BodyParser(&req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = Validate.Struct(req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = r.service.ProcessSubscriptionDeleteRequest(c.UserContext(), &req.UserId)
	if err != nil {
		return errors.Wrap(err, "failed to process subscription delete request")
	}

	zerolog.Ctx(c.UserContext()).
		Info().
		Str("subscription_request", string(c.Body())).
		Msg("new.delete.subscription.request")

	return c.JSON(fiber.Map{"status": "success"})
}

func (r *Presentation) updateSubscription(c *fiber.Ctx) error {
	var req *services.SubscriptionRequest

	err := c.BodyParser(&req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = Validate.Struct(req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = r.service.ProcessSubscriptionUpdateRequest(c.UserContext(), req)
	if err != nil {
		return errors.Wrap(err, "failed to update subscription")
	}

	zerolog.Ctx(c.UserContext()).
		Info().
		Str("subscription_request", string(c.Body())).
		Msg("new.update.subscription.request")

	return c.JSON(fiber.Map{"status": "success"})
}
