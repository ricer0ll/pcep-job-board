package utils

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
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
