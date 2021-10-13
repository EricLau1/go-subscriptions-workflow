package webtokens

import (
	"github.com/gofiber/fiber/v2"
	"go-subscriptions-workflow/security/tokens"
)

func GetToken(ctx *fiber.Ctx) (*tokens.TokenPayload, error) {
	token := ctx.Request().Header.Peek("Authorization")
	return tokens.Parse(string(token))
}
