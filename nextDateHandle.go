package main

import (
	"net/http"
	"time"

	"github.com/ogrestudies/go_final_project/task"
)

func nextDateHandle(res http.ResponseWriter, req *http.Request) {
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")
	now, _ := time.Parse(task.DateLayout, req.FormValue("now"))

	nextDate, _ := task.NextDate(now, date, repeat)
	res.Write([]byte(nextDate))

}
