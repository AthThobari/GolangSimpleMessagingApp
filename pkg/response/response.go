package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	SuccessMessage = "success"
)

func SendsuccessResponse(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(Response{Message: SuccessMessage, Data: data})
}

func SendFailureResponse(ctx *fiber.Ctx,httpCode int, message string, data interface{}) error {
	return ctx.Status(httpCode).JSON(Response{Message: message, Data: data})
}