package api

import (
	"final-project/internal/planner"
	"net/http"
	"time"
)

const dateFormat = "20060102"

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repStr := r.FormValue("repeat")

	var now time.Time
	var err error
	var result string

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	result, err = planner.NextDate(now, dateStr, repStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}
