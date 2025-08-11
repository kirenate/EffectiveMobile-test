package presentations

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"main.go/services"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())

func (r *Presentation) postSubscription(c *fiber.Ctx) error {
	var req *services.SubscriptionRequest

	err := c.BodyParser(&req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = Validate.StructExcept(req, "UserId")
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = r.service.ProcessSubscriptionRequest(c.UserContext(), req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: err.Error()}
	}

	log.Info().
		Interface("subscription_request", req).
		Msg("new.subscription.request")

	return c.JSON(fiber.Map{"status": "success"})
}

func (r *Presentation) getSubscription(c *fiber.Ctx) error {
	var req *services.SubscriptionRequest

	err := c.BodyParser(&req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = Validate.StructPartial(req, "UserId")
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	sub, err := r.service.ProcessSubscriptionGetRequest(c.UserContext(), req.UserId)
	if err != nil {
		return errors.Wrap(err, "failed to process subscription get request")
	}

	log.Info().
		Interface("user_id", req.UserId).
		Msg("new.get.subscription.request")

	return c.JSON(sub)
}

func (r *Presentation) deleteSubscription(c *fiber.Ctx) error {
	var req *services.SubscriptionRequest

	err := c.BodyParser(&req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = Validate.StructPartial(req, "UserId")
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = r.service.ProcessSubscriptionDeleteRequest(c.UserContext(), &req.UserId)
	if err != nil {
		return errors.Wrap(err, "failed to process subscription delete request")
	}

	log.Info().
		Interface("user_id", req.UserId).
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

	log.Info().
		Interface("subscription_request", req).
		Msg("new.update.subscription.request")

	return c.JSON(fiber.Map{"status": "success"})
}

func (r *Presentation) listSubscriptions(c *fiber.Ctx) error {
	sub, err := r.service.ProcessSubscriptionListRequest(c.UserContext())
	if err != nil {
		return errors.Wrap(err, "failed to process subscription get request")
	}

	log.Info().Msg("new.list.subscriptions.request")

	return c.JSON(sub)
}

func (r *Presentation) subscriptionCostUserId(c *fiber.Ctx) error {
	var req *services.SubscriptionRequest

	err := c.BodyParser(&req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = Validate.StructPartial(req, "UserId")
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	cost, err := r.service.ProcessSubscriptionCostUserId(c.UserContext(), req.UserId)

	return c.JSON(fiber.Map{"cost": cost})
}

func (r *Presentation) subscriptionCostServiceName(c *fiber.Ctx) error {
	var req *services.SubscriptionRequest

	err := c.BodyParser(&req)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	err = Validate.StructPartial(req, "ServiceName")
	if err != nil {
		return &fiber.Error{Code: fiber.StatusUnprocessableEntity, Message: err.Error()}
	}

	cost, err := r.service.ProcessSubscriptionCostServiceName(c.UserContext(), req.ServiceName)

	return c.JSON(fiber.Map{"cost": cost})
}
