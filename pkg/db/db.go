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

// Init - функция для инициализации базы данных.
// Она выполняет следующие шаги:
// 1. Проверка наличия файла базы данных с использованием функции os.Stat. Если файл не существует, то устанавливается флаг install в true.
// 2. Открытие соединения с базой данных с использованием функции sql.Open. Если открытие соединения не удалось, то возвращается ошибка.
// 3. Проверка соединения с базой данных с использованием функции DB.Ping. Если проверка не удалась, то возвращается ошибка.
// 4. Если флаг install установлен в true, то выполняется установка базы данных с использованием функции DB.Exec. Если установка не удалась, то возвращается ошибка.
// 5. Возвращение ошибки, если она возникла при инициализации базы данных.
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

// Close - функция для закрытия соединения с базой данных.
// Она выполняет следующие шаги:
// 1. Проверка наличия соединения с базой данных. Если соединение отсутствует, то возвращается nil.
// 2. Закрытие соединения с базой данных с использованием функции DB.Close. Если закрытие соединения не удалось, то возвращается ошибка.
// 3. Возвращение ошибки, если она возникла при закрытии соединения с базой данных.
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
