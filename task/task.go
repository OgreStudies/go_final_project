package task

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"     //формат записи времени
const SearchLayout = "02.01.2006" //формат фремени для запроса

// Функция возвращает следующую дату события
// now — время от которого ищется ближайшая дата;
// date — исходное время в формате 20060102, от которого начинается отсчёт повторений;
// repeat — правило повторения
//   - Если правило не указано, отмеченная выполненной задача будет удаляться из таблицы;
//   - d <число> — задача переносится на указанное число дней. Максимально допустимое число равно 400. Примеры:
//     d 1 — каждый день;
//     d 7 — для вычисления следующей даты добавляем семь дней;
//     d 60 — переносим на 60 дней.
//   - y — задача выполняется ежегодно. Этот параметр не требует дополнительных уточнений. При выполнении задачи дата перенесётся на год вперёд.
//   - w <через запятую от 1 до 7> — задача назначается в указанные дни недели, где 1 — понедельник, 7 — воскресенье. Например:
//     w 7 — задача перенесётся на ближайшее воскресенье;
//     w 1,4,5 — задача перенесётся на ближайший понедельник, четверг или пятницу;
//     w 2,3 — задача перенесётся на ближайший вторник или среду.
//   - m <через запятую от 1 до 31,-1,-2> [через запятую от 1 до 12] — задача назначается в указанные дни месяца. При этом вторая последовательность чисел опциональна и указывает на определённые месяцы. Например:
//     m 4 — задача назначается на четвёртое число каждого месяца;
//     m 1,15,25 — задача назначается на 1-е, 15-е и 25-е число каждого месяца;
//     m -1 — задача назначается на последний день месяца;
//     m -2 — задача назначается на предпоследний день месяца;
//     m 3 1,3,6 — задача назначается на 3-е число января, марта и июня;
//     m 1,-1 2,8 — задача назначается на 1-е и последнее число число февраля и авгуcта.
func NextDate(now time.Time, date string, repeat string) (string, error) {
	//Парсим строку базовой даты
	startDate, err := time.Parse(DateLayout, date)
	if err != nil {
		return "", fmt.Errorf("ошибка формата исходного времени: %w", err)
	}
	//Разбиваем правило повтора на составляющие
	ruleParts := strings.Split(repeat, " ")

	//Формирование результата на основании правила повтора
	switch ruleParts[0] {
	case "d": //Через определённое количество дней
		//Если количество частей правила повторения < 2 - Ошибка формата правила повторений
		if len(ruleParts) < 2 {
			return "", fmt.Errorf("ошибка формата правила повторения: %q", repeat)
		}
		nDays, err := strconv.ParseInt(ruleParts[1], 10, 64)
		if err != nil || nDays < 1 || nDays > 400 {
			return "", fmt.Errorf("ошибка формата правила повторения: %q", repeat)
		}

		//Поиск следующей даты
		nextDate := startDate.AddDate(0, 0, int(nDays))
		for {

			if nextDate.After(now) {
				return nextDate.Format(DateLayout), nil
			}
			nextDate = nextDate.AddDate(0, 0, int(nDays))
		}
	case "y": //Раз в год
		//Поиск следующей даты
		nextDate := startDate.AddDate(1, 0, 0)
		for {
			if nextDate.After(now) || nextDate.Equal(now) {
				return nextDate.Format(DateLayout), nil
			}
			nextDate = nextDate.AddDate(1, 0, 0)
		}

	case "w": //задача назначается в указанные дни недели
		//Если количество частей правила повторения < 2 - Ошибка формата правила повторений
		if len(ruleParts) < 2 {
			return "", fmt.Errorf("ошибка формата правила повторения: %q", repeat)
		}
		//Парсим дни недели в которые нужно повторение
		weekDays := [7]int{0, 0, 0, 0, 0, 0, 0}
		firstDay := int(0)
		for i, d := range strings.Split(ruleParts[1], ",") {
			val, err := strconv.ParseInt(d, 10, 64)
			if err != nil || val < 1 || val > 7 {
				return "", fmt.Errorf("ошибка формата правила повторения: %q", repeat)
			}
			//запомнить первый разрешённый день
			if i == 0 {
				firstDay = int(val)
			}
			weekDays[val-1] = 1
		}

		//первый разрешённый день не может быть равным 0
		if firstDay == 0 {
			return "", fmt.Errorf("ошибка формата правила повторения: %q", repeat)
		}

		nextDate := now
		if startDate.After(now) {
			nextDate = startDate
		}

		//Дни недели от 1 до 7, 1 = понеднльник
		wdNow := int(nextDate.Weekday())
		if wdNow == 0 {
			wdNow = 7
		}
		nextDateOffset := int(0)

		//Поиск дней недели после текущего
		for i := wdNow; i < len(weekDays); i++ {
			if weekDays[i] == 1 {
				nextDateOffset = i + 1 - wdNow
				break
			}
		}
		//Если не нашлось дней повтора на текущей неделе, то переходим на следующую неделю и назначаем первый из доступных дней
		if nextDateOffset == 0 {
			nextDateOffset = firstDay + (7 - wdNow)
		}

		nextDate = nextDate.AddDate(0, 0, nextDateOffset)
		return nextDate.Format(DateLayout), nil
	case "m": //задача назначается в указанные дни месяца
		//Если количество частей правила повторения < 2 - Ошибка формата правила повторений
		if len(ruleParts) < 2 {
			return "", fmt.Errorf("ошибка формата правила повторения: %q", repeat)
		}
		//Парсим месяцы в которые нужно повторение
		firstMonth := int(0)
		yearMonths := [12]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

		nextDate := now
		if startDate.After(now) {
			nextDate = startDate
		}

		if len(ruleParts) > 2 {
			for i, m := range strings.Split(ruleParts[2], ",") {
				val, err := strconv.ParseInt(m, 10, 64)
				if err != nil || val < 1 || val > 12 {
					return "", fmt.Errorf("ошибка формата правила повторения: %q", repeat)
				}
				if i == 0 {
					firstMonth = int(val)
				}
				yearMonths[val-1] = 1
			}
		}
		//Если нет ограничений по месяцам, то доступны все месяцы
		if firstMonth == 0 {
			yearMonths = [12]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
			firstMonth = 1
		}

		//Парсим дни месяца в которые нужно повторение
		monthDays := [12][31]int{}
		for _, d := range strings.Split(ruleParts[1], ",") {
			val, err := strconv.ParseInt(d, 10, 64)
			if err != nil || val < -2 || val > 31 || val == 0 {
				return "", fmt.Errorf("ошибка формата правила повторения: %q", repeat)
			}

			for m, em := range yearMonths { //Дни отмечаются только в разрешённых месяцах
				if em == 1 {
					if val > 0 {
						//Если значение дня > 0, то отмечаем дни с начала месяца
						switch time.Month(m + 1) {
						case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
							monthDays[m][val-1] = 1
						case time.April, time.June, time.September, time.November:
							if val <= 30 {
								monthDays[m][val-1] = 1
							}
						case time.February:
							if val <= 29 {
								monthDays[m][val-1] = 1
							}
						}
					} else {
						//Если значение дня < 0, то отмечаем дни с конца месяца
						switch time.Month(m + 1) {
						case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
							monthDays[m][31+val] = 1
						case time.April, time.June, time.September, time.November:
							monthDays[m][30+val] = 1
						case time.February:
							//Для февраля отмечаем как для месяца с кол-вом дней = 31
							monthDays[m][31+val] = 1
						}
					}
				}
			}
		}

		//Функция определяет является ли год високосным
		isLeap := func(year int) bool {
			if year%400 == 0 {
				return true
			} else if year%100 == 0 {
				return true
			} else if year%4 == 0 {
				return true
			}
			return false
		}

		//Функция фозвращает первый разрешённый день больше dNow
		//Или -1, если день не найден
		getDay := func(dNow int, mNow time.Month, leapYear bool, daysMask [31]int) int {
			retVal := -1
			for i := dNow; i < 31; i++ {
				if daysMask[i] == 1 {
					//Если февраль, то различаем високосный/не високосный год
					if mNow == time.February {
						if !leapYear && ((i + 1) == 29) { //Для не високосного года пропускаем 29-е число
							continue
						} else if (i + 1) < 30 { //Если i соответствует реальному дню возвращаем этот день
							retVal = i + 1
						} else { //Если i не соответствует реальному дню то выполняем смещение номера дня на основании информации о том является ли год високосным
							if leapYear {
								retVal = i - 1
							} else {
								retVal = i - 2
							}
						}
						break
					} else { //для остальных месяцев возвращаем день
						retVal = i + 1
						break
					}
				}
			}
			return retVal

		}

		//Проверяем есть ли текущий месяц в перечне разрешённых
		if yearMonths[nextDate.Month()-1] == 1 {
			//Ищем разрешённые дни после текущего
			retVal := getDay(nextDate.Day(), time.Month(nextDate.Month()), isLeap(nextDate.Year()), monthDays[nextDate.Month()-1])
			//Если нашли такой день, то вернуть эту дату
			if retVal > 0 {
				return time.Date(nextDate.Year(), nextDate.Month(), retVal, 0, 0, 0, 0, time.UTC).Format(DateLayout), nil
			}
		}

		//Поиск следующего разрешённого месяца
		for i := nextDate.Month(); i < 12; i++ {
			if yearMonths[i] == 1 {
				//Ищем первый разрешённый день
				retVal := getDay(0, time.Month(nextDate.Month()), isLeap(nextDate.Year()), monthDays[i])
				//Если нашли такой день, то вернуть эту дату
				if retVal > 0 {
					return time.Date(nextDate.Year(), i+1, retVal, 0, 0, 0, 0, time.UTC).Format(DateLayout), nil
				}
			}
		}

		//Первый разрешённый день следующего года
		retVal := getDay(0, time.Month(firstMonth), isLeap(nextDate.Year()+1), monthDays[firstMonth-1])
		if retVal > 0 {
			return time.Date(nextDate.Year()+1, time.Month(firstMonth), retVal, 0, 0, 0, 0, time.UTC).Format(DateLayout), nil
		}

	}

	return "", fmt.Errorf("ошибка формата правила повторения: %q", repeat)
}

// Структура - описание задачи
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Валидация и корректировка полей Task
func (t *Task) TaskFieldCheckAndCorrect() (Task, error) {
	newTask := Task{}
	newTask.ID = t.ID
	newTask.Title = t.Title
	newTask.Date = t.Date
	newTask.Comment = t.Comment
	newTask.Repeat = t.Repeat

	//Поле Title не может быть пустым
	if newTask.Title == "" {
		return newTask, fmt.Errorf("поле 'title' не может быть пустым")
	}
	//Если Date пустое - вмето него подставляется сегодняшнее число
	if newTask.Date == "" {
		newTask.Date = time.Now().Format(DateLayout)
	} else {
		//Если Date не пустое - Поле Date должно иметь формат "20260102
		_, err := time.Parse(DateLayout, newTask.Date)
		if err != nil {
			return newTask, fmt.Errorf("ошибка формата поля 'date'")
		}
	}

	if newTask.Repeat == "" { //Пустое правило повторения
		if newTask.Date < time.Now().Format(DateLayout) { //Если task.Date меньше текщей даты, то task.Date устанавливается сегодняшним числом
			newTask.Date = time.Now().Format(DateLayout)
		}
	} else { //Не Пустое правило повторения
		nextDate, err := NextDate(time.Now(), newTask.Date, newTask.Repeat)
		if err != nil { //Ошибка правила повторений

			return newTask, fmt.Errorf("ошибка формата поля 'repeat'")
		}
		if newTask.Date < time.Now().Format(DateLayout) { //Если task.Date меньше текщей даты, то task.Date устанавливается nextDate
			newTask.Date = nextDate
		}
	}
	return newTask, nil
}

// Структура - описание хранилища задач
type TaskStorage struct {
	db *sql.DB
}

// Возвращает TaskStorage с заданным в `db` хранилищем
func NewTaskstorage(db *sql.DB) TaskStorage {
	return TaskStorage{db: db}
}

// Добавляет задачу в хранилище
func (ts TaskStorage) AddTask(task Task) (int64, error) {

	//Проверка и коррекция задачи
	checkedTask, err := task.TaskFieldCheckAndCorrect()
	if err != nil {
		return 0, err
	}

	//Добавление новой задачи в базу
	dbres, err := ts.db.Exec("INSERT INTO scheduler(date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", checkedTask.Date),
		sql.Named("title", checkedTask.Title),
		sql.Named("comment", checkedTask.Comment),
		sql.Named("repeat", checkedTask.Repeat),
	)

	if err != nil {
		return 0, fmt.Errorf("ошибка добавления записи: %w", err)

	}

	//Идентификатор добвленной задачи
	id, _ := dbres.LastInsertId()
	return id, nil

}

func (ts TaskStorage) GetTaskById(id int64) (Task, error) {
	sqlString := "SELECT * FROM scheduler WHERE id = :id;"
	pRow := ts.db.QueryRow(sqlString, sql.Named("id", id))
	task := Task{}
	err := pRow.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, fmt.Errorf("ошибка получения записи с id: %d: %w", id, err)
	}
	return task, nil

}

// Обновляет задачу в хранилище
func (ts TaskStorage) UpdateTask(id int64, taskData *Task) (Task, error) {
	//проверка наличия записи с требуемым id
	storedtask, err := ts.GetTaskById(id)
	if err != nil {
		return storedtask, fmt.Errorf("запись с id: %v не найдена", id)
	}

	//Проверка и коррекция данных для обновления
	checkedTask, err := taskData.TaskFieldCheckAndCorrect()
	if err != nil {
		return checkedTask, fmt.Errorf("запись с id: %v ошибка формата данных для обновления: %w", id, err)
	}
	//Добавление идентификатора
	checkedTask.ID = fmt.Sprintf("%d", id)

	//Если запись найдена - обновить её
	sqlString := "UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat  WHERE id = :id"
	_, err = ts.db.Exec(sqlString,
		sql.Named("id", id),
		sql.Named("date", checkedTask.Date),
		sql.Named("title", checkedTask.Title),
		sql.Named("comment", checkedTask.Comment),
		sql.Named("repeat", checkedTask.Repeat),
	)

	if err != nil {

		return checkedTask, err
	}

	return checkedTask, nil
}

// Удалаяет задачу из хранилища
func (ts TaskStorage) DeleteTask(id int64) error {
	//проверка наличия записи с требуемым id
	_, err := ts.GetTaskById(id)
	if err != nil {
		return fmt.Errorf("запись с id: %v не найдена", id)
	}

	//Удаление записи если она найдена
	sqlString := "DELETE FROM scheduler WHERE id = :id;"
	_, err = ts.db.Exec(sqlString, sql.Named("id", id))
	if err != nil {

		return fmt.Errorf("ошибка удаления записи с id: %v : %w", id, err)
	}

	return nil
}

// Получает ближайшие задачи. Количество ограничено параметром num
func (ts TaskStorage) GetLastTasks(num int, searchReq string) ([]Task, error) {

	tasks := make([]Task, 0, num)

	var pRows *sql.Rows
	var err error

	if searchReq != "" { //Если поисковый запрос
		sqlString := ""
		t, parceErr := time.Parse(SearchLayout, searchReq)
		if parceErr != nil { //Если не удалось преобразовать запрос в дату, то поисковый запрос по Title
			sqlString = "SELECT * FROM scheduler WHERE title LIKE :searchReq ORDER BY date ASC LIMIT :maxlim;"
		} else { //Иначе поисковый запрос по дате
			searchReq = t.Format(DateLayout) //Преобразование в формат базы данных
			sqlString = "SELECT * FROM scheduler WHERE date LIKE :searchReq ORDER BY date ASC LIMIT :maxlim;"
		}
		searchReq = "%" + searchReq + "%"
		pRows, err = ts.db.Query(sqlString,
			sql.Named("searchReq", searchReq),
			sql.Named("maxlim", num),
		)
	} else { //Не поисковый запрос
		pRows, err = ts.db.Query("SELECT * FROM scheduler ORDER BY date ASC LIMIT :maxlim;",
			sql.Named("maxlim", num),
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
