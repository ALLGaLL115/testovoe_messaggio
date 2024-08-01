package message

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ALLGaLL115/testovoe-messaggio/internal/domain/dto"
	"github.com/ALLGaLL115/testovoe-messaggio/internal/domain/models"
	"github.com/ALLGaLL115/testovoe-messaggio/internal/handlers"
	"github.com/ALLGaLL115/testovoe-messaggio/internal/lib/logger/sl"

	// "github.com/ALLGaLL115/testovoe-messaggio/internal/lib/middleware"

	"github.com/go-chi/chi/v5/middleware"
)

type MessageHandler struct {
	log            *slog.Logger
	messageService MessageService
	AppID          int32
}

type MessageService interface {
	Create(ctx context.Context, message dto.Message) (int64, error)
	GetById(ctx context.Context, messageID int64) (models.Message, error)
	UpdateByID(ctx context.Context, message dto.Message) (int64, error)
	DeleteByID(ctx context.Context, messageID int64) (int64, error)
}

func New(log *slog.Logger, messageService MessageService, AppID int32) *MessageHandler {
	return &MessageHandler{
		log:            log,
		messageService: messageService,
		AppID:          AppID,
	}
}

func (h *MessageHandler) Create(ctx context.Context) http.HandlerFunc {
	const op = "handlers/message/Create"

	return func(w http.ResponseWriter, r *http.Request) {
		h.log = h.log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var message dto.MessageRequest
		err := json.NewDecoder(r.Body).Decode(&message)

		if err != nil {
			h.log.Error("failed to decode request body", sl.Err(err))
			handlers.ErrorResponse(w, r, 400, "bad request")
		}

		if err := message.Validate(); err != nil {
			h.log.Error("failed to validate message", sl.Err(err))
			handlers.ErrorResponse(w, r, 422, err.Error())
		}

		messageModel := dto.Message{
			Text:      message.Text,
			CreatedAt: time.Now().UTC(),
		}

		messageID, err := h.messageService.Create(ctx, messageModel)

		if err != nil {
			h.log.Error("failed to create message", sl.Err(err))
			handlers.ErrorResponse(w, r, 500, "failed to create message")
			return
		}

		handlers.SuccessRespnose(w, r, 200, map[string]any{
			"message":    "message successfuly created",
			"meesage_id": messageID,
		})

	}
}

// func (h *MessageHandler) GetByID(ctx context.Context) http.HandlerFunc{
// 	const op = "handlers/message/GetByID"

// 	return func(w http.ResponseWriter, r *http.Request) {
// 		h.log = h.log.With(
// 			slog.String("op", op),
// 			slog.String("request_id", middleware.GetReqID(r.Context())),
// 		)

// 	}
// }
