package handlers

import "github.com/ogrestudies/go_final_project/internal/tasks"

var todoStorage *tasks.TaskStorage

func SetStorage(storage *tasks.TaskStorage) {
	todoStorage = storage
}
