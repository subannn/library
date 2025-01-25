package externalApiClient

import (
	"log/slog"
)

type apiClient struct {
	logger *slog.Logger
}

func NewAPIClient(logger *slog.Logger) *apiClient {
	return &apiClient{
		logger: logger,
	}
}

func (apiClient *apiClient) GetSongDetails(group, songName string) (releaseDate, text, link string, err error) {
	return "01.01.2004", "hello_world", "something.com", nil
}
