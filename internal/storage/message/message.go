package message

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ALLGaLL115/testovoe-messaggio/internal/domain/dto"
	"github.com/ALLGaLL115/testovoe-messaggio/internal/domain/models"
	"github.com/ALLGaLL115/testovoe-messaggio/internal/lib/logger/sl"
	"github.com/ALLGaLL115/testovoe-messaggio/internal/lib/storage/query"
	"github.com/jackc/pgx/v5"
)

type MessageDB struct {
	log *slog.Logger
}

func newMessageDB(log *slog.Logger) *MessageDB {
	return &MessageDB{
		log: log,
	}
}

const (
	messageTable = "messages"
)

var (
	ErrMessageNotFound  = errors.New("message not found")
	ErrMessagesNotFound = errors.New("messages not found")
)

func (m *MessageDB) CreateMessage(ctx context.Context, tx pgx.Tx, message dto.Message) (int64, error) {
	const op = "storage.message.CreateMessage"

	q := fmt.Sprintf(`
	INSERT INTO %s
		(text, created_at)
	VALUES 
		($1, $2,)
			RETURNING id;
	`, messageTable)

	m.log.Debug("create message query", slog.String("query", query.QueryToString(q)))

	var messageID int64

	err := tx.QueryRow(ctx, q, message.Text, message.CreatedAt).Scan(&messageID)

	if err != nil {
		m.log.Error("failed to create message", sl.OpError(op, err))
		return 0, err
	}

	return messageID, nil
}

func (m *MessageDB) GetMessageById(ctx context.Context, tx pgx.Tx, messageID int64) (models.Message, error) {
	const op = "storage.message.CreateMessage"

	q := fmt.Sprintf(`
	SELECT 
		id, text, created_at
	FROM %s
	WHERE id = $1
	`, messageTable)

	m.log.Debug("get message by id query:", slog.String("query", query.QueryToString(q)))

	var message models.Message
	err := tx.QueryRow(ctx, q, messageID).Scan(&message.ID, &message.Text, &message.CreatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Message{}, ErrMessageNotFound
		}
		m.log.Error("failed to get message by id", sl.OpError(op, err))
		return models.Message{}, err
	}

	return message, nil

}

func (m *MessageDB) UpdateMessageByID(ctx context.Context, tx pgx.Tx, ID uint64, message dto.Message) (int64, error) {
	const op = "storage.message.UpdateMessageByID"

	q := fmt.Sprintf(`
	UPDATE %s
	SET  text=$2
	WHERE id = $1
	RETURNING id

	`, messageTable)

	m.log.Debug("update message by id: ", slog.String("query", query.QueryToString(q)))

	var messageID int64
	err := tx.QueryRow(ctx, q, ID, message.Text).Scan(&messageID)

	if err != nil {
		m.log.Error("failed to update a message: ", sl.OpError(op, err))
		return 0, err
	}

	return messageID, nil
}

func (m *MessageDB) DeleteMessageByID(ctx context.Context, tx pgx.Tx, ID int64) (int64, error) {
	const op = "storage.message.DeleteMessageByID"

	q := fmt.Sprintf(`
	DELETE FROM %s
	WHERE id = $1
	RETURNING id
	`, messageTable)

	m.log.Debug("delete message by id: ", slog.String("query", query.QueryToString(q)))

	var messageId int64
	err := tx.QueryRow(ctx, q, ID).Scan(&messageId)

	if err != nil {
		m.log.Error("failed to delete message: ", sl.OpError(op, err))
		return 0, err
	}

	return messageId, nil

}
