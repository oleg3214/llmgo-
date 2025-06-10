package http

import (
	"context"
	"log"
	"net/http"

	"github.com/Ablyamitov/userbot-core/internal/domain/dto/request"
	"github.com/Ablyamitov/userbot-core/internal/domain/dto/response/wrapper"
	"github.com/Ablyamitov/userbot-core/internal/usecase/userbot"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	usecase *userbot.Usecase
}

func NewHandler(u *userbot.Usecase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Post("/api/chat/messages", h.GetMessages)
	app.Post("/api/chat/send-message", h.SendMessage)
}

func (h *Handler) GetMessages(c *fiber.Ctx) error {
	method := "UserbotHTTPServer.GetMessagesFromChannel"
	var req *request.MessagesRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("%s: Error parsing request body: %v", method, err)
		return c.Status(http.StatusBadRequest).JSON(wrapper.ResponseWrapper{
			Code:  http.StatusBadRequest,
			Error: wrapper.Error{Message: "incorrect data"},
		})
	}
	massages, err := h.usecase.GetMessages(context.Background(), req.Username, int(req.Limit))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(wrapper.ResponseWrapper{
			Code:  http.StatusInternalServerError,
			Error: wrapper.Error{Message: "failed to get messages"},
		})
	}

	var result []*wrapper.Message
	for _, m := range massages {

		/**/

		result = append(result, wrapper.FromEntity(&m))
		/**/
		//result = append(result, &wrapper.Message{
		//	MessageId: int64(m.ID),
		//	Text:      m.Text,
		//	Date:      int64(m.Date),
		//})
	}

	return c.Status(http.StatusOK).JSON(wrapper.ResponseWrapper{
		Status: true,
		Code:   http.StatusOK,
		Data:   wrapper.MessagesResponse{Messages: result}, //map to dto
		Error:  nil,
	})

}

func (h *Handler) SendMessage(c *fiber.Ctx) error {
	method := "UserbotHTTPServer.SendMessage"
	var req *request.SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("%s: Error parsing request body: %v", method, err)
		return c.Status(http.StatusBadRequest).JSON(wrapper.ResponseWrapper{
			Code:  http.StatusBadRequest,
			Error: wrapper.Error{Message: "incorrect data"},
		})
	}
	err := h.usecase.SendMessage(context.Background(), req.Username, req.Text)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(wrapper.ResponseWrapper{
			Code:  http.StatusInternalServerError,
			Error: wrapper.Error{Message: "send error"},
		})
	}
	return c.Status(http.StatusOK).JSON(wrapper.ResponseWrapper{
		Status: true,
		Code:   http.StatusOK,
	})

}
