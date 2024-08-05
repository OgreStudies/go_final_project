package tasks

import (
	"fmt"
	"time"
)

const DateLayout = "20060102"     //формат записи времени
const SearchLayout = "02.01.2006" //формат фремени для запроса

// Структура - описание задачи
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

// Валидация и корректировка полей Task
func (t *Task) TaskFieldCheckAndCorrect() (Task, error) {
	newTask := *t

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
