package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/copyleftdev/snaptrack/pkg/snapshot"
)

type DBInterface interface {
	GetLastSnapshot(url string) (snapshot.Snapshot, error)
	InsertSnapshot(s snapshot.Snapshot) error
}

type DB struct {
	conn *sql.DB
}

func InitDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	query := `
    CREATE TABLE IF NOT EXISTS snapshots (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT NOT NULL,
        hash TEXT NOT NULL,
        html TEXT NOT NULL,
        status_code INTEGER,
        request_headers TEXT,
        response_headers TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err = conn.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (d *DB) GetLastSnapshot(url string) (snapshot.Snapshot, error) {
	row := d.conn.QueryRow(`
        SELECT id, url, hash, html, status_code, request_headers, response_headers, created_at
          FROM snapshots
         WHERE url = ?
      ORDER BY created_at DESC
         LIMIT 1
    `, url)

	var s snapshot.Snapshot
	var reqHeadersJSON, respHeadersJSON []byte
	var createdStr string

	err := row.Scan(
		&s.ID,
		&s.URL,
		&s.Hash,
		&s.HTML,
		&s.StatusCode,
		&reqHeadersJSON,
		&respHeadersJSON,
		&createdStr,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return s, errors.New("no rows")
		}
		return s, err
	}

	if len(reqHeadersJSON) > 0 {
		_ = json.Unmarshal(reqHeadersJSON, &s.RequestHeaders)
	}
	if len(respHeadersJSON) > 0 {
		_ = json.Unmarshal(respHeadersJSON, &s.ResponseHeaders)
	}

	t, _ := time.Parse("2006-01-02 15:04:05", createdStr)
	s.CreatedAt = t
	return s, nil
}

func (d *DB) InsertSnapshot(s snapshot.Snapshot) error {
	reqHeadersJSON, _ := json.Marshal(s.RequestHeaders)
	respHeadersJSON, _ := json.Marshal(s.ResponseHeaders)

	_, err := d.conn.Exec(`
        INSERT INTO snapshots (url, hash, html, status_code, request_headers, response_headers, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `,
		s.URL,
		s.Hash,
		s.HTML,
		s.StatusCode,
		string(reqHeadersJSON),
		string(respHeadersJSON),
		s.CreatedAt.Format("2006-01-02 15:04:05"),
	)
	return err
}

func (d *DB) Close() error {
	return d.conn.Close()
}
