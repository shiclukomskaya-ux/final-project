package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"os"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(32) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
	);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
`

var DB *sql.DB

func Init(dbFile string) error {
	_, errStat := os.Stat(dbFile)
	var install bool
	if errStat != nil {
		install = true
	}

	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
