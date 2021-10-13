package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go-subscriptions-workflow/api/webtokens"
	"go-subscriptions-workflow/services/users/service"
	"go-subscriptions-workflow/types"
	"go-subscriptions-workflow/util"
	"net/http"
)

type usersHandlers struct {
	usersService   service.UsersService
	inputValidator *validator.Validate
}

func RegisterUsersHandlers(usersService service.UsersService, app *fiber.App) {
	h := &usersHandlers{usersService: usersService, inputValidator: validator.New()}
	app.Post("/signup", h.PostSignUp)
	app.Post("/login", h.PostLogin)
	app.Get("/users/:id", h.GetUser)
	app.Put("/users/:id/credit", h.PutCredit)
}

func (h *usersHandlers) PostSignUp(ctx *fiber.Ctx) error {
	in := new(types.CreateUserInput)
	err := ctx.BodyParser(in)
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	err = h.inputValidator.Struct(in)
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	in.Email = util.NormalizeEmail(in.Email)
	out, err := h.usersService.CreateUser(ctx.Context(), in)
	if err != nil {
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.
		Status(http.StatusCreated).
		JSON(out)
}

func (h *usersHandlers) PostLogin(ctx *fiber.Ctx) error {
	in := new(types.LoginInput)
	err := ctx.BodyParser(in)
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	err = h.inputValidator.Struct(in)
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	in.Email = util.NormalizeEmail(in.Email)
	out, err := h.usersService.Login(ctx.Context(), in)
	if err != nil {
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.
		Status(http.StatusOK).
		JSON(out)
}

func (h *usersHandlers) GetUser(ctx *fiber.Ctx) error {
	out, err := h.usersService.GetUser(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.
		Status(http.StatusOK).
		JSON(out)
}

func (h *usersHandlers) PutCredit(ctx *fiber.Ctx) error {
	token, err := webtokens.GetToken(ctx)
	if err != nil {
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(fiber.Map{"error": err.Error()})
	}
	if token.UserID != ctx.Params("id") {
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(fiber.Map{"error": "invalid request"})
	}
	in := new(types.CreditInput)
	err = ctx.BodyParser(in)
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	err = h.inputValidator.Struct(in)
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	in.UserID = token.UserID
	out, err := h.usersService.Credit(ctx.Context(), in)
	if err != nil {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.
		Status(http.StatusOK).
		JSON(out)
}
