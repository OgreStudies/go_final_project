package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// Открывает или создаёт базу данных для хранения задач
func OpenStorage(storagePath string) (*sql.DB, error) {

	//Проверка наличия файла по указанному пути
	_, err := os.Stat(storagePath)

	//Если файл не найден, то создать новую базу данных
	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных %q: %w", storagePath, err)
	}

	//Инициализация новой базы данных
	if install {
		_, err := db.Exec(
			`CREATE TABLE scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date VARCHAR(8) NOT NULL DEFAULT "",
				title VARCHAR(256) NOT NULL DEFAULT "",
				comment TEXT NOT NULL DEFAULT "",
				repeat VARCHAR(128) NOT NULL DEFAULT ""
			);  
		`)
		if err != nil {
			db.Close()
			os.Remove(storagePath)
			return nil, fmt.Errorf("ошибка создания таблицы 'scheduler': %w", err)
		}
		_, err = db.Exec(`CREATE INDEX scheduler_date ON scheduler (date);`)
		if err != nil {
			db.Close()
			os.Remove(storagePath)
			return nil, fmt.Errorf("ошибка создания индекса по полю 'date': %w", err)
		}
	}
	//База данных открыта успешно, вернуть указатель на соединение
	return db, err
}
