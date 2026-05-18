package event

import (
	"log/slog"

	"github.com/disgoorg/disgo/events"
)

func OnReady(event *events.Ready) {
	slog.Info("Bot connected and ready!")
}
