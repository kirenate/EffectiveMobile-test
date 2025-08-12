package presentations

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"main.go/services"
)

type Presentation struct {
	service *services.Service
}

func NewPresentation(service *services.Service) *Presentation {
	return &Presentation{service: service}
}

func (r *Presentation) BuildApp() *fiber.App {
	app := fiber.New(fiber.Config{})
	app.Use(recover2.New(recover2.Config{EnableStackTrace: true}))

	app.Post("/subscription", r.postSubscription)
	app.Get("/subscription", r.getSubscription)
	app.Delete("/subscription", r.deleteSubscription)
	app.Put("/subscription", r.updateSubscription)

	app.Get("/subscriptions", r.listSubscriptions)

	app.Get("/subscription/cost/user_id", r.subscriptionCostUserId)
	app.Get("/subscription/cost/service_name", r.subscriptionCostServiceName)

	return app
}
