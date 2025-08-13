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

	app.Get("/openapi.yaml", r.openapi)
	app.Get("/docs", r.swagger)

	app.Post("/subscriptions", r.postSubscription)
	app.Get("/subscriptions", r.getSubscription)
	app.Delete("/subscriptions", r.deleteSubscription)
	app.Put("/subscriptions", r.updateSubscription)

	app.Get("/subscription-list", r.listSubscriptions)

	app.Get("/subscription/cost/user-id", r.subscriptionCostUserId)
	app.Get("/subscription/cost/service-name", r.subscriptionCostServiceName)

	return app
}
