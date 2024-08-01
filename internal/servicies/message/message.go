package message

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ALLGaLL115/testovoe-messaggio/internal/domain/dto"
	"github.com/ALLGaLL115/testovoe-messaggio/internal/domain/models"
	"github.com/ALLGaLL115/testovoe-messaggio/internal/lib/logger/sl"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MessageService struct {
	log        *slog.Logger
	messagesDB MessagesDB
	pool       *pgxpool.Pool
}

type MessagesDB interface {
	CreateMessage(ctx context.Context, tx pgx.Tx, message dto.Message) (int64, error)
	GetMessageById(ctx context.Context, tx pgx.Tx, messageID int64) (models.Message, error)
	UpdateMessageByID(ctx context.Context, tx pgx.Tx, message dto.Message) (int64, error)
	DeleteMessageByID(ctx context.Context, tx pgx.Tx, messageID int64) (int64, error)
}

func NewMessageService(log *slog.Logger, messagesDB MessagesDB, pool *pgxpool.Pool) *MessageService {
	return &MessageService{
		log:        log,
		messagesDB: messagesDB,
		pool:       pool,
	}
}

func (service *MessageService) Create(ctx context.Context, message dto.Message) (int64, error) {
	const op = "servicies.message.CreateMessage"

	tx, err := service.pool.Begin(ctx)

	if err != nil {
		service.log.Error("failed to start transaction", sl.OpError(op, err))
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		if err := tx.Commit(ctx); err != nil {
			service.log.Error("failed to commit transaction", sl.OpError(op, err))
			return
		}
	}()

	messageID, err := service.messagesDB.CreateMessage(ctx, tx, message)
	if err != nil {
		service.log.Error("failed to create message", sl.OpError(op, err))
		return 0, err
	}

	if messageID == 0 {
		err = errors.New("message is empty")
		service.log.Error("failed to create message", sl.OpError(op, err))
		return 0, err
	}

	return messageID, nil

}

func (service *MessageService) GetById(ctx context.Context, tx pgx.Tx, messageID int64) (models.Message, error) {
	const op = "servicies.message.GetMessageById"
	tx, err := service.pool.Begin(ctx)

	if err != nil {
		service.log.Error("failed to start transaction", sl.OpError(op, err))
		return models.Message{}, err
	}

	defer tx.Rollback(ctx)

	message, err := service.messagesDB.GetMessageById(ctx, tx, messageID)

	if err != nil {
		service.log.Error("failed get message by id", sl.OpError(op, err))
		return models.Message{}, err
	}

	if message == (models.Message{}) {
		err := errors.New("message is empty")
		service.log.Error("failed get message by id", sl.OpError(op, err))
		return models.Message{}, err
	}

	return message, nil

}

func (service *MessageService) UpdateByID(ctx context.Context, tx pgx.Tx, message dto.Message) (int64, error) {
	const op = "servicies.message.UpdateMessageByID"

	tx, err := service.pool.Begin(ctx)

	if err != nil {
		service.log.Error("failed to start transaction", sl.OpError(op, err))
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}

		if err := tx.Commit(ctx); err != nil {
			service.log.Error("failed to start transaction", sl.OpError(op, err))
			return
		}
	}()

	messageID, err := service.messagesDB.UpdateMessageByID(ctx, tx, message)

	if err != nil {
		service.log.Error("failed update message", sl.OpError(op, err))
		return 0, err
	}
	if messageID == 0 {
		err := errors.New("message is empty")
		service.log.Error("failed update message", sl.OpError(op, err))
		return 0, err
	}

	return messageID, nil
}

func (service *MessageService) DeleteByID(ctx context.Context, tx pgx.Tx, messageID int64) (int64, error) {
	const op = "servicies.message.DeleteMessageByID"

	tx, err := service.pool.Begin(ctx)

	if err != nil {
		service.log.Error("failed to delete message", sl.OpError(op, err))
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}

		if err := tx.Commit(ctx); err != nil {
			service.log.Error("failed to commit transaction", sl.OpError(op, err))
			return
		}
	}()

	messageID, err = service.messagesDB.DeleteMessageByID(ctx, tx, messageID)

	if err != nil {
		service.log.Error("failed to delete message", sl.OpError(op, err))
	}

	if messageID == 0 {
		err := errors.New("message empty")
		service.log.Error("failed to delete message", sl.OpError(op, err))
		return 0, err
	}

	return messageID, nil

}
