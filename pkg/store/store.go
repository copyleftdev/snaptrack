package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/copyleftdev/snaptrack/pkg/snapshot"
)

// DBInterface is the interface for storing and retrieving snapshots.
type DBInterface interface {
	GetLastSnapshot(url string) (snapshot.Snapshot, error)
	InsertSnapshot(s snapshot.Snapshot) error
	// Optionally: GetDistinctURLs(), GetSnapshotsForURL(), etc.
}

// DB is our concrete implementation holding a *sql.DB connection.
type DB struct {
	conn *sql.DB
}

// InitDB opens (or creates) a SQLite DB at dbPath, migrates tables, returns *DB.
func InitDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	query := `
    CREATE TABLE IF NOT EXISTS snapshots (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT NOT NULL,
        hash TEXT NOT NULL,
        html TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	if _, err := conn.Exec(query); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &DB{conn: conn}, nil
}

// GetLastSnapshot returns the most recent snapshot for a URL.
func (d *DB) GetLastSnapshot(url string) (snapshot.Snapshot, error) {
	row := d.conn.QueryRow(`
        SELECT id, url, hash, html, created_at
        FROM snapshots
        WHERE url = ?
        ORDER BY created_at DESC
        LIMIT 1
    `, url)

	var s snapshot.Snapshot
	var createdStr string
	err := row.Scan(&s.ID, &s.URL, &s.Hash, &s.HTML, &createdStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return s, errors.New("no rows")
		}
		return s, err
	}

	// parse time
	t, _ := time.Parse("2006-01-02 15:04:05", createdStr)
	s.CreatedAt = t
	return s, nil
}

// InsertSnapshot inserts a new snapshot into the DB.
func (d *DB) InsertSnapshot(s snapshot.Snapshot) error {
	_, err := d.conn.Exec(`
        INSERT INTO snapshots (url, hash, html, created_at)
        VALUES (?, ?, ?, ?)
    `, s.URL, s.Hash, s.HTML, s.CreatedAt.Format("2006-01-02 15:04:05"))
	return err
}

// (Optional) Close the DB connection.
func (d *DB) Close() error {
	return d.conn.Close()
}
