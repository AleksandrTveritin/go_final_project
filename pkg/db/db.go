package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	if install {
		if _, err = DB.Exec(schema); err != nil {
			return err
		}
	}
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
