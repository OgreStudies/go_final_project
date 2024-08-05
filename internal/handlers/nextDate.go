package handlers

import (
	"net/http"
	"time"

	"github.com/ogrestudies/go_final_project/internal/tasks"
)

func NextDateHandle(res http.ResponseWriter, req *http.Request) {

	date := req.FormValue("date")
	repeat := req.FormValue("repeat")
	now, _ := time.Parse(tasks.DateLayout, req.FormValue("now"))

	nextDate, _ := tasks.NextDate(now, date, repeat)
	res.Write([]byte(nextDate))

}
