package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/ALLGaLL115/testovoe-messaggio/internal/validator"
)

type Message struct {
	// ID        int64     `json:"id" validate:"required"`
	Text      string    `json:"text" validate:"required"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
}

func (m *Message) Validate() error {
	m.Text = strings.TrimSpace(m.Text)

	if err := validator.Validate(m); err != "" {
		return fmt.Errorf("validation error: %s", err)
	}
	return nil

}

type MessageRequest struct {
	Text string `json:"text" validate:"required"`
}

func (messageRequest *MessageRequest) Validate() error {
	messageRequest.Text = strings.TrimSpace(messageRequest.Text)

	if err := validator.Validate(messageRequest); err != "" {
		return fmt.Errorf("validation error: %s", err)
	}
	return nil
}
