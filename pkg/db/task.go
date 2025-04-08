package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/AleksandrTveritin/go_final_project/pkg/config"
)

type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask - функция для добавления задачи в базу данных.
// Она выполняет следующие шаги:
// 1. Проверка корректности задачи с использованием функции validateTask. Если задача некорректна, то возвращается ошибка.
// 2. Выполнение SQL-запроса для добавления задачи в базу данных с использованием функции DB.Exec. Если выполнение запроса не удалось, то возвращается ошибка.
// 3. Возвращение идентификатора добавленной задачи с использованием функции res.LastInsertId. Если получение идентификатора не удалось, то возвращается ошибка.
// 4. Возвращение ошибки, если она возникла при добавлении задачи в базу данных.
func AddTask(task *Task) (int64, error) {
	if err := validateTask(task); err != nil {
		return 0, err
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) 
              VALUES (?, ?, ?, ?)`

	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("failed to add task: %w", err)
	}

	return res.LastInsertId()
}

func GetTasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat 
              FROM scheduler 
              WHERE date >= ?`

	params := []interface{}{time.Now().Format(config.DateFormat)}

	query += ` ORDER BY date ASC, id ASC LIMIT ?`
	params = append(params, limit)

	rows, err := DB.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return tasks, nil
}

// GetTasks - функция для получения списка задач из базы данных.
// Она выполняет следующие шаги:
// 1. Формирование SQL-запроса для получения задач из базы данных. Если указан параметр search, то добавляется условие поиска по заголовку и комментарию задачи.
// 2. Выполнение SQL-запроса с использованием функции DB.Query. Если выполнение запроса не удалось, то возвращается ошибка.
// 3. Чтение результатов запроса с использованием функции rows.Next. Если чтение результатов не удалось, то возвращается ошибка.
// 4. Преобразование результатов запроса в структуры Task с использованием функции rows.Scan. Если преобразование не удалось, то возвращается ошибка.
// 5. Возвращение списка задач и ошибки, если она возникла при получении задач из базы данных.
func GetTask(id int64) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat 
              FROM scheduler 
              WHERE id = ?`

	var task Task
	err := DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// UpdateTask - функция для обновления задачи в базе данных.
// Она выполняет следующие шаги:
// 1. Проверка корректности задачи с использованием функции validateTask. Если задача некорректна, то возвращается ошибка.
// 2. Формирование SQL-запроса для обновления задачи в базе данных.
// 3. Выполнение SQL-запроса с использованием функции DB.Exec. Если выполнение запроса не удалось, то возвращается ошибка.
// 4. Получение количества обновленных строк с использованием функции res.RowsAffected. Если получение количества строк не удалось, то возвращается ошибка.
// 5. Если количество обновленных строк равно нулю, то возвращается ошибка "task not found".
// 6. Возвращение ошибки, если она возникла при обновлении задачи в базе данных.
func UpdateTask(task *Task) error {
	if err := validateTask(task); err != nil {
		return err
	}

	query := `UPDATE scheduler 
              SET date = ?, title = ?, comment = ?, repeat = ?
              WHERE id = ?`

	res, err := DB.Exec(query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// validateTask - функция для проверки корректности задачи.
// Она выполняет следующие шаги:
// 1. Проверка наличия заголовка задачи. Если заголовок отсутствует, то возвращается ошибка.
// 2. Проверка формата даты задачи. Если формат даты некорректен, то возвращается ошибка.
// 3. Возвращение ошибки, если задача некорректна.
func validateTask(task *Task) error {
	if task.Title == "" {
		return fmt.Errorf("title is required")
	}

	if _, err := time.Parse(config.DateFormat, task.Date); err != nil {
		return fmt.Errorf("invalid date format")
	}

	return nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id int64) error {
	res, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}

// UpdateDate обновляет только дату задачи
func UpdateDate(id int64, date string) error {
	_, err := DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", date, id)
	if err != nil {
		return fmt.Errorf("failed to update task date: %v", err)
	}
	return nil
}
