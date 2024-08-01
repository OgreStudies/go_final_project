package handlers

import "github.com/ogrestudies/go_final_project/internal/taskstorage"

var todoStorage *taskstorage.TaskStorage

func SetStorage(storage *taskstorage.TaskStorage) {
	todoStorage = storage
}
