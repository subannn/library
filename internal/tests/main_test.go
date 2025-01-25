package tests

import (
	"EffectiveMobileTestTask/internal/config"
	"EffectiveMobileTestTask/internal/externalApiClient"
	"EffectiveMobileTestTask/internal/handlers"
	"EffectiveMobileTestTask/internal/libaryDB"
	"EffectiveMobileTestTask/internal/models"
	"EffectiveMobileTestTask/internal/server"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHandlers(t *testing.T) {
	logger := NewLogger()
	ctx, cancel := context.WithCancel(context.Background())
	go getHanlder(ctx)
	time.Sleep(10 * time.Second)

	data := models.Song{
		SongName: "testingName",
		Group:    "testingGroup",
	}
	jsonData, err := json.Marshal(data)
	resp, err := http.Post("http://localhost:8085/saveSong", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var songID int64
	err = json.Unmarshal(body, &songID)
	if err != nil {
		t.Fatalf("Ошибка десериализации JSON: %v", err)
	}
	strSongId := strconv.FormatInt(songID, 10)
	logger.Info("Song saved", "songID", strSongId)

	// update
	data1 := models.SongQuery{
		Group:       "testNewgroup",
		SongName:    "testNewName",
		ReleaseDate: "2024-12-05",
		Text:        "lirycs1/lirycs2/lirycs3/lirycs4/lirycs5/lirycs6/lirycs7",
		Link:        "https://www.lirycsTestLink.org/",
	}
	jsonData, err = json.Marshal(data1)

	url := "http://localhost:8085/updateSong/" + strSongId

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Ошибка создания PUT-запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка выполнения PUT-запроса: %v", err)
	}
	defer resp.Body.Close()

	resp, err = http.Get("http://localhost:8085/getSongText/" + strSongId)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)

	var responseText []string
	err = json.Unmarshal(body, &responseText)

	if !reflect.DeepEqual(responseText, strings.Split(data1.Text, "/")) {
		logger.Info("RESSPPONSE", "TEXT=", responseText)
		t.Fatal(responseText)
	}
	cancel()

}

func getHanlder(ctx context.Context) {
	logger := NewLogger()

	logger.Info("Logger initialized")

	cfg := config.NewConfig()
	cfg.DBConfig = &config.DBConfig{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "1",
		Database: "musicLibary",
		SSLMode:  "disable",
	}
	logger.Info("Configuration initialized")

	logger.Info("Waiting for postgres container...")
	time.Sleep(time.Second * 5)

	db := libaryDB.NewDB(logger, cfg)
	logger.Info("DB initialized")

	apiClient := externalApiClient.NewAPIClient(logger)
	logger.Info("External API client initialized")

	handler := handlers.NewHandler(db, apiClient, logger)
	logger.Info("Handler initialized")

	srv := server.NewServer(logger, handler)
	logger.Info("Server initialized")

	go srv.Start(cfg.ServerConfig.Port)
	time.Sleep(time.Second * 5)

	<-ctx.Done()
	logger.Info("Test context canceled")
}

func NewLogger() *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(handler)
}
