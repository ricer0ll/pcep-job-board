package utils

import (
	"log/slog"
	"os"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/joho/godotenv"
)

var (
	channelID string
	loaded    bool
)

func GetBotConfig() bot.ConfigOpt {
	configOpts := bot.WithGatewayConfigOpts(
		gateway.WithIntents(
			gateway.IntentGuilds,
			gateway.IntentGuildMessages,
		),
	)

	return configOpts
}

func GetDiscordChannelID() string {
	if loaded {
		return channelID
	}

	err := godotenv.Load()
	if err != nil {
		slog.Warn("Unable to load .env file. Fetching system environment variables...")
	}

	channelID = os.Getenv("CHANNEL_ID")
	if channelID == "" {
		panic("Unable to load CHANNEL_ID environment variable.")
	}
	loaded = true
	return channelID
}
