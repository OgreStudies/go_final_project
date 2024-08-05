package tasks

import (
	"database/sql"
	"fmt"
	"os"
	"time"

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

// Структура - описание хранилища задач
type TaskStorage struct {
	db              *sql.DB
	ReturnTaskLimit int
}

// Возвращает TaskStorage с заданным в `db` хранилищем
func NewTaskstorage(db *sql.DB, retuntTaskLimit int) TaskStorage {
	return TaskStorage{db: db, ReturnTaskLimit: retuntTaskLimit}
}

// Добавляет задачу в хранилище
func (ts TaskStorage) AddTask(task *Task) (int64, error) {

	if task == nil {
		return 0, fmt.Errorf("указатель task не может быть nil")
	}

	//Добавление новой задачи в базу
	dbres, err := ts.db.Exec("INSERT INTO scheduler(date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)

	if err != nil {
		return 0, fmt.Errorf("ошибка добавления записи: %w", err)
	}

	//Идентификатор добвленной задачи
	id, _ := dbres.LastInsertId()
	return id, nil

}

func (ts TaskStorage) GetTaskById(id int64) (Task, error) {
	sqlString := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id;"
	pRow := ts.db.QueryRow(sqlString, sql.Named("id", id))
	task := Task{}
	err := pRow.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, fmt.Errorf("ошибка получения записи с id: %d: %w", id, err)
	}
	return task, nil

}

// Обновляет задачу в хранилище
func (ts TaskStorage) UpdateTask(id int64, taskData *Task) error {

	if taskData == nil {
		return fmt.Errorf("указатель taskData не может быть nil")
	}

	//Добавление идентификатора
	taskData.ID = fmt.Sprintf("%d", id)

	//Если запись найдена - обновить её
	sqlString := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat  WHERE id = :id"
	res, err := ts.db.Exec(sqlString,
		sql.Named("id", id),
		sql.Named("date", taskData.Date),
		sql.Named("title", taskData.Title),
		sql.Named("comment", taskData.Comment),
		sql.Named("repeat", taskData.Repeat),
	)

	if err != nil {

		return err
	}

	if nrows, err := res.RowsAffected(); nrows == 0 && err == nil {
		return fmt.Errorf("запись с id: %v не найдена", id)
	}

	return nil
}

// Удалаяет задачу из хранилища
func (ts TaskStorage) DeleteTask(id int64) error {
	//Удаление записи если она найдена
	sqlString := "DELETE FROM scheduler WHERE id = :id;"
	res, err := ts.db.Exec(sqlString, sql.Named("id", id))
	if err != nil {

		return fmt.Errorf("ошибка удаления записи с id: %v : %w", id, err)
	}

	if nrows, err := res.RowsAffected(); nrows == 0 && err == nil {
		return fmt.Errorf("запись с id: %v не найдена", id)
	}

	return nil
}

// Получает ближайшие задачи. Количество ограничено параметром num
func (ts TaskStorage) GetLastTasks(searchReq string) ([]Task, error) {

	tasks := make([]Task, 0, ts.ReturnTaskLimit)

	var pRows *sql.Rows
	var err error

	if searchReq != "" { //Если поисковый запрос
		sqlString := ""
		t, parceErr := time.Parse(SearchLayout, searchReq)
		if parceErr != nil { //Если не удалось преобразовать запрос в дату, то поисковый запрос по Title
			sqlString = "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :searchReq OR comment LIKE :searchReq ORDER BY date ASC LIMIT :maxlim;"
		} else { //Иначе поисковый запрос по дате
			searchReq = t.Format(DateLayout) //Преобразование в формат базы данных
			sqlString = "SELECT id, date, title, comment, repeat FROM scheduler WHERE date LIKE :searchReq ORDER BY date ASC LIMIT :maxlim;"
		}
		searchReq = "%" + searchReq + "%"
		pRows, err = ts.db.Query(sqlString,
			sql.Named("searchReq", searchReq),
			sql.Named("maxlim", ts.ReturnTaskLimit),
		)
	} else { //Не поисковый запрос
		pRows, err = ts.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT :maxlim;",
			sql.Named("maxlim", ts.ReturnTaskLimit),
		)
	}
	if err != nil {
		return tasks, err
	}

	defer pRows.Close()
	// заполняем срез Tasks данными из таблицы
	for pRows.Next() {
		var task Task
		err := pRows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, task)
	}

	if err := pRows.Err(); err != nil {
		return tasks, err
	}
	return tasks, nil
}
