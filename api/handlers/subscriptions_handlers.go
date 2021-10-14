package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go-subscriptions-workflow/api/webtokens"
	"go-subscriptions-workflow/rmq"
	"go-subscriptions-workflow/services/subscriptions/service"
	"go-subscriptions-workflow/services/subscriptions/shared"
	"go-subscriptions-workflow/types"
	"net/http"
)

type subscriptionsHandlers struct {
	subsClient service.SubscriptionsClient
	producer   rmq.Producer
}

func RegisterSubscriptionsHandlers(subsClient service.SubscriptionsClient, producer rmq.Producer, app *fiber.App) {
	h := &subscriptionsHandlers{
		subsClient: subsClient,
		producer:   producer,
	}
	app.Post("/subscriptions", h.PostStartSubscription)
	app.Put("/subscriptions/:id/cancel", h.PutCancelSubscription)
	app.Get("/subscriptions", h.GetSubscriptions)
	app.Get("/subscriptions/:id", h.GetSubscription)
}

func (h *subscriptionsHandlers) PostStartSubscription(ctx *fiber.Ctx) error {
	token, err := webtokens.GetToken(ctx)
	if err != nil {
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(fiber.Map{"error": err.Error()})
	}
	req := &types.StartSubscriptionRequest{UserID: token.UserID}
	options := &rmq.PublisherOptions{
		ExchangeName: shared.ExchangeName,
		Persistent:   true,
	}
	err = h.producer.Send(options, rmq.NewMessage(req))
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.
		Status(http.StatusAccepted).
		JSON(fiber.Map{"status": "sent"})
}

func (h *subscriptionsHandlers) PutCancelSubscription(ctx *fiber.Ctx) error {
	token, err := webtokens.GetToken(ctx)
	if err != nil {
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(fiber.Map{"error": err.Error()})
	}
	req := &types.CancelSubscriptionRequest{
		ID:     ctx.Params("id"),
		UserID: token.UserID,
	}
	options := &rmq.PublisherOptions{
		ExchangeName: shared.ExchangeName,
		Persistent:   true,
	}
	err = h.producer.Send(options, rmq.NewMessage(req))
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.
		Status(http.StatusAccepted).
		JSON(fiber.Map{"status": "sent"})
}

func (h *subscriptionsHandlers) GetSubscriptions(ctx *fiber.Ctx) error {
	out, err := h.subsClient.GetSubscriptions(ctx.Context())
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.
		Status(http.StatusOK).
		JSON(out)
}

func (h *subscriptionsHandlers) GetSubscription(ctx *fiber.Ctx) error {
	req := &types.GetSubscriptionRequest{ID: ctx.Params("id")}
	out, err := h.subsClient.GetSubscription(ctx.Context(), req)
	if err != nil {
		return ctx.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.
		Status(http.StatusOK).
		JSON(out)
}

