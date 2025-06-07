package persistence

import (
	sql "database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func SqliteConnection() (*sql.DB, error) {
	return sql.Open("sqlite3", "./app.db")
}

func Migrate(con *sql.DB) {
	con.Exec(`CREATE TABLE IF NOT EXISTS device (	
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id VARCHAR(6) UNIQUE,
		device_url VARCHAR(20) UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
}
