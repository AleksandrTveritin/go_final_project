package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   int
	DBFile string
}

var AppConfig *Config

// LoadConfig - функция для загрузки конфигурации приложения.
// Она выполняет следующие шаги:
// 1. Загрузка переменных окружения из файла .env с использованием функции godotenv.Load.
// 2. Создание структуры Config с значениями по умолчанию.
// 3. Проверка наличия переменной окружения TODO_PORT. Если она указана, то ее значение преобразуется в целочисленное значение и используется в качестве порта приложения. Если преобразование не удалось, то выводится сообщение об ошибке и используется значение по умолчанию.
// 4. Проверка наличия переменной окружения TODO_DBFILE. Если она указана, то ее значение используется в качестве имени файла базы данных.
// 5. Возвращение ошибки, если она возникла при загрузке конфигурации.
func LoadConfig() error {
	_ = godotenv.Load(".env")

	AppConfig = &Config{
		Port:   7540,           // Значение по умолчанию
		DBFile: "scheduler.db", // Значение по умолчанию
	}

	if portStr := os.Getenv("TODO_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Printf("Invalid port value: %s, using default %d", portStr, AppConfig.Port)
		} else {
			AppConfig.Port = port
		}
	}

	if dbFile := os.Getenv("TODO_DBFILE"); dbFile != "" {
		AppConfig.DBFile = dbFile
	}

	return nil
}
