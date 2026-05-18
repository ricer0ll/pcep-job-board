package utils

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/joho/godotenv"
)

func GetDiscordToken() (*string, error) {
	var token string

	// .env file
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Unable to load .env file. Fetching from GCP secret manager.")
	} else {
		token = os.Getenv("DISCORD_TOKEN")
		if len(token) != 0 {
			return &token, nil
		}
	}

	// GCP secret manager
	token, err = getGcpSecret()
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func getGcpSecret() (string, error) {
	name := "projects/274174031721/secrets/discord-token/versions/1"

	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("Faild to create secretmanager client: %w", err)
	}
	defer client.Close()

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("Failed to get secret: %w", err)
	}

	return string(result.Payload.Data), nil
}
