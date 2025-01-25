package libaryDB

import "database/sql"

func (db *DB) insertGroupAndSongWithTransaction(groupName, songName, releaseDate, lyrics, link string) (int64, error) {
	// start of transaction
	tx, err := db.db.Beginx()
	if err != nil {
		db.logger.Info("Failed to start transaction: %v", err)
		return 0, err
	}

	var groupID int64
	queryInsertGroup := `INSERT INTO groups (name) VALUES ($1) RETURNING id`
	err = tx.Get(&groupID, queryInsertGroup, groupName)
	if err != nil {
		db.logger.Debug("Failed to insert group: %v", err)
		_ = tx.Rollback() // Откат транзакции в случае ошибки
		return 0, err
	}
	db.logger.Info("Group inserted with ID: %d", groupID)

	// Insert song using groupId
	var songID int64
	queryInsertSong := `
		INSERT INTO songs (group_id, song_name, release_date, lyrics, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	err = tx.Get(&songID, queryInsertSong, groupID, songName, releaseDate, lyrics, link)
	if err != nil {
		db.logger.Debug("Failed to insert song: %v", err)
		_ = tx.Rollback() // Откат транзакции в случае ошибки
		return 0, err
	}
	db.logger.Info("Song inserted with ID: %d", songID)

	if err := tx.Commit(); err != nil {
		db.logger.Debug("Failed to commit transaction: %v", err)
		return 0, err
	}

	return songID, nil
}

func (db *DB) inserOnlySong(groupID int64, songName, releaseDate, text, link string) (int64, error) {
	queryInsertSong := `
		INSERT INTO songs (group_id, song_name, release_date, lyrics, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	var songID int64
	err := db.db.Get(&songID, queryInsertSong, groupID, songName, releaseDate, text, link)
	if err != nil {
		db.logger.Debug("Error inserting song: %v", err)
		return 0, err
	}

	db.logger.Info("Song saved successfully")
	return songID, nil
}

func (db *DB) insertOnlyGroup(groupName string) (int64, error) {
	groupExists, groupId, err := db.ifGroupExists(groupName)
	if err != nil {
		return 0, err
	}
	if !groupExists {
		var groupID int64
		queryInsertGroup := `INSERT INTO groups (name) VALUES ($1) RETURNING id`
		err = db.db.Get(&groupID, queryInsertGroup, groupName)
		if err != nil {
			db.logger.Debug("Error inserting group: %v", err)
			return 0, err
		}
		db.logger.Info("Group inserted with ID: %d", groupID)
		return groupID, nil
	}

	return groupId, nil
}

func (db *DB) ifGroupExists(groupName string) (bool, int64, error) {
	var groupID int64

	queryCheckGroup := `SELECT id FROM groups WHERE name = $1`
	err := db.db.Get(&groupID, queryCheckGroup, groupName)
	if err == sql.ErrNoRows {
		return false, -1, nil
	} else if err != nil {
		return false, -1, err
	}

	return true, groupID, nil
}
