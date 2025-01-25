package libaryDB

import (
	"EffectiveMobileTestTask/internal/models"
	"fmt"
	"strings"
)

func (db *DB) SaveSong(group, songName, releaseDate, text, link string) (int64, error) {
	groupExists, groupID, err := db.ifGroupExists(group)
	if err != nil {
		db.logger.Debug("Error while checking group", "group", group, "error", err)
		return -1, err
	}
	var songId int64
	if groupExists {
		songId, err = db.inserOnlySong(groupID, songName, releaseDate, text, link)
		if err != nil {
			db.logger.Debug("Error while saving song")
			return -1, err
		}
	} else {
		songId, err = db.insertGroupAndSongWithTransaction(group, songName, releaseDate, text, link)
		if err != nil {
			db.logger.Debug("Error while saving song")
			return -1, err
		}
	}

	return songId, nil
}

func (db *DB) UpdateSong(songId int64, group, songName, releaseDate, lyrics, link string) error {
	var groupID int64
	var err error

	if group != "" {
		groupID, err = db.insertOnlyGroup(group)
		if err != nil {
			db.logger.Debug("Error while saving song", "group", group, "error", err)
			return err
		}
	}

	updateQuery := `UPDATE songs SET `
	setClauses := []string{}
	params := []interface{}{}
	paramIndex := 1

	if songName != "" {
		setClauses = append(setClauses, fmt.Sprintf("song_name = $%d", paramIndex))
		params = append(params, songName)
		paramIndex++
	}
	if releaseDate != "" {
		setClauses = append(setClauses, fmt.Sprintf("release_date = $%d", paramIndex))
		params = append(params, releaseDate)
		paramIndex++
	}
	if lyrics != "" {
		setClauses = append(setClauses, fmt.Sprintf("lyrics = $%d", paramIndex))
		params = append(params, lyrics)
		paramIndex++
	}
	if link != "" {
		setClauses = append(setClauses, fmt.Sprintf("link = $%d", paramIndex))
		params = append(params, link)
		paramIndex++
	}
	if group != "" {
		setClauses = append(setClauses, fmt.Sprintf("group_id = $%d", paramIndex))
		params = append(params, groupID)
		paramIndex++
	}

	if len(setClauses) == 0 {
		db.logger.Warn("No fields to update", "songId", songId)
		return nil
	}

	updateQuery += strings.Join(setClauses, ", ")
	updateQuery += fmt.Sprintf(" WHERE id = $%d", paramIndex)
	params = append(params, songId)

	db.logger.Info("Executing query", "query", updateQuery, "params", params)

	_, err = db.db.Exec(updateQuery, params...)
	if err != nil {
		db.logger.Debug("Failed to execute update query", "error", err)
		return err
	}

	return nil
}

func (db *DB) DeleteSong(songId int64) error {
	var exists bool
	queryCheck := "SELECT EXISTS (SELECT 1 FROM songs WHERE id=$1)"
	err := db.db.QueryRow(queryCheck, songId).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check song existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("song with id %d does not exist", songId)
	}

	queryDeleteSong := "DELETE FROM songs WHERE id=$1"
	res, err := db.db.Exec(queryDeleteSong, songId)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}
	db.logger.Info("Deleted song", "id", songId, "results", res)

	return nil
}

func (db *DB) GetSongText(songId int64) ([]string, error) {
	queryCheck := "SELECT lyrics FROM songs WHERE id=$1"
	var text string
	err := db.db.Get(&text, queryCheck, songId)
	if err != nil {
		db.logger.Debug("Error while getting song text", "id", songId, "error", err)
		return nil, err
	}
	parts := strings.Split(text, "/")

	return parts, err

}

func (db *DB) GetSongs(group, songName, releaseDate string, limit, offset int) (songs []models.SongResponse, err error) {
	query := `
		SELECT 
			s.id AS song_id, 
			g.id AS group_id, 
			s.song_name, 
			g.name AS group_name, 
			s.release_date
		FROM 
			songs s
		JOIN 
			groups g ON s.group_id = g.id
		WHERE 
			(COALESCE(:group, '') = '' OR g.name ILIKE :group) AND
			(COALESCE(:song_name, '') = '' OR s.song_name ILIKE :song_name) AND
			(COALESCE(:release_date, '') = '' OR s.release_date = :release_date)
		ORDER BY 
			s.release_date DESC
		LIMIT :limit OFFSET :offset`

	params := map[string]interface{}{
		"group":        filterValue(group),
		"song_name":    filterValue(songName),
		"release_date": filterValue(releaseDate),
		"limit":        limit,
		"offset":       offset,
	}

	rows, err := db.db.NamedQuery(query, params)
	if err != nil {
		db.logger.Debug("Error while getting songs", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song models.SongResponse
		if err := rows.StructScan(&song); err != nil {
			db.logger.Debug("Error while scanning song", "error", err)
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func filterValue(value string) interface{} {
	if value == "" {
		return nil
	}
	return "%" + value + "%"
}
