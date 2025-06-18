// internal/storage/sqlite.go
package storage

import (
	"database/sql"
	"fmt"
	"webscraper/internal/scraper"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

const dbPath = "scraper.db" // Database file path

// InitDB initializes the SQLite database and creates the links table if it doesn't exist.
func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Check if the table exists, and create it if not.
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			text TEXT,
			href TEXT NOT NULL
		)
	`)
	if err != nil {
		db.Close() // Close the database if table creation fails
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}

// InsertLink inserts a single link into the database.
func InsertLink(db *sql.DB, url string, link scraper.Link) error {
	_, err := db.Exec("INSERT INTO links (url, text, href) VALUES (?, ?, ?)",
		url, link.Text, link.Href)
	if err != nil {
		return fmt.Errorf("failed to insert link: %w", err)
	}
	return nil
}

// InsertLinks inserts multiple links into the database in a single transaction.
// This is more efficient than inserting them one by one.
func InsertLinks(db *sql.DB, url string, links []scraper.Link) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if any insertion fails

	stmt, err := tx.Prepare("INSERT INTO links (url, text, href) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, link := range links {
		_, err = stmt.Exec(url, link.Text, link.Href)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
