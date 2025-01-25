package handlers

import (
	"EffectiveMobileTestTask/internal/models"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"

	_ "github.com/swaggo/echo-swagger"
)

type Service interface {
	SaveSong(group, songName, releaseDate, text, link string) (int64, error)
	UpdateSong(songId int64, group, songName, releaseDate, text, link string) error
	DeleteSong(songId int64) error
	GetSongText(songId int64) ([]string, error)
	GetSongs(group, songName, releaseDate string, limit, offset int) (songs []models.SongResponse, err error) // may be will change return variables
}

type ExternalAPIClient interface {
	GetSongDetails(group, songName string) (releaseDate, text, link string, err error)
}

type Handler struct {
	service           Service
	externalAPIClient ExternalAPIClient
	logger            *slog.Logger
}

func NewHandler(service Service, externalAPIClient ExternalAPIClient, logger *slog.Logger) *Handler {
	return &Handler{
		service:           service,
		externalAPIClient: externalAPIClient,
		logger:            logger,
	}
}

// SaveSong сохраняет песню в базу данных
// @Summary Сохранить песню
// @Description Сохраняет информацию о песне, включая дату выпуска, текст и ссылку на внешние данные
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Данные песни"
// @Success 200 {object} map[string]int "ID сохраненной песни"
// @Failure 400 {object} map[string]string "Ошибка привязки данных"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /saveSong [post]
func (h *Handler) SaveSong(c echo.Context) error {
	var requestSong models.Song
	if err := c.Bind(&requestSong); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	releaseDate, text, link, err := h.externalAPIClient.GetSongDetails(requestSong.Group, requestSong.SongName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	songId, err := h.service.SaveSong(requestSong.Group, requestSong.SongName, releaseDate, text, link)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	h.logger.Info("Song saved")
	return c.JSON(http.StatusOK, songId)
}

// UpdateSong обновляет информацию о песне
// @Summary Обновить песню
// @Description Обновляет информацию о песне по ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body models.SongQuery true "Данные для обновления песни"
// @Success 200 {string} string "Успешное обновление"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /updateSong/{id} [put]
func (h *Handler) UpdateSong(c echo.Context) error {
	songId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	var requestSongUpdate models.SongQuery
	if err := c.Bind(&requestSongUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = h.service.UpdateSong(int64(songId), requestSongUpdate.Group, requestSongUpdate.SongName, requestSongUpdate.ReleaseDate, requestSongUpdate.Text, requestSongUpdate.Link)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	h.logger.Info("Song updated")
	return c.JSON(http.StatusOK, nil)
}

// DeleteSong удаляет песню
// @Summary Удалить песню
// @Description Удаляет песню по ID
// @Tags songs
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {string} string "Успешное удаление"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /deleteSong/{id} [delete]
func (h *Handler) DeleteSong(c echo.Context) error {
	songId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = h.service.DeleteSong(songId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	h.logger.Info("Song deleted")
	return c.JSON(http.StatusOK, nil)
}

// GetSongText возвращает текст песни
// @Summary Получить текст песни
// @Description Возвращает текст песни по ID
// @Tags songs
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {string} string "Текст песни"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /getSongText/{id} [get]
func (h *Handler) GetSongText(c echo.Context) error {
	songId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	text, err := h.service.GetSongText(songId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if text == nil {
		return c.JSON(http.StatusOK, "No text found")
	}

	h.logger.Info("Song's text extracted")
	return c.JSON(http.StatusOK, text)
}

// GetSongs возвращает список песен
// @Summary Получить список песен
// @Description Возвращает список песен с возможностью фильтрации и пагинации
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Название группы"
// @Param song_name query string false "Название песни"
// @Param release_date query string false "Дата выпуска"
// @Param limit query int false "Количество записей" default(10)
// @Param offset query int false "Смещение для пагинации" default(0)
// @Success 200 {array} models.Song "Список песен"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /getSongs [get]
func (h *Handler) GetSongs(c echo.Context) error {
	var songsParamsRequest models.SongQuery
	if err := c.Bind(&songsParamsRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	limitParam := c.QueryParam("limit")
	offsetParam := c.QueryParam("offset")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}
	songs, err := h.service.GetSongs(songsParamsRequest.Group, songsParamsRequest.SongName, songsParamsRequest.ReleaseDate, limit, offset)
	if err != nil {
		h.logger.Info("something went wrong", "err", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	h.logger.Info("Songs retrieved")
	return c.JSON(http.StatusOK, songs)
}
