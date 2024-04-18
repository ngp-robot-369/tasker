package server

import (
	"encoding/json"
	"net/http"
	"tasker/tasks"
)

type (
	HttpHandler func(http.ResponseWriter, *http.Request)
	Route       struct {
		Name    string
		Pattern string
		Handler HttpHandler
	}
)

func GetRoutes() []Route {
	return []Route{
		{
			"Creates new task to perform request",
			"POST /task",
			serveCreateTask,
		},
		{
			"Returns task status if any found",
			"GET /task/{id}",
			serveGetTaskStatus,
		},
	}
}

func serveCreateTask(w http.ResponseWriter, r *http.Request) {
	var request tasks.TaskRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"id": tasks.NewTask(request),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func serveGetTaskStatus(w http.ResponseWriter, r *http.Request) {
	var (
		id     = r.PathValue("id")
		status = tasks.GetStatus(id)
	)

	if status == nil {
		http.Error(w, "Task id not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
