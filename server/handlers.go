package server

import (
	"encoding/json"
	"net/http"
	"tasker/tasks"
)

func (h Handlers) serveCreateTask(w http.ResponseWriter, r *http.Request) {
	var request tasks.TaskRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.Tasker.NewTask(request)
	if err != nil {
		http.Error(w, "Create task failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id": id,
	})
}

func (h Handlers) serveGetTaskStatus(w http.ResponseWriter, r *http.Request) {
	var (
		id          = r.PathValue("id")
		status, err = h.Tasker.GetStatus(id)
	)
	if err != nil {
		http.Error(w, "Get task failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if status == nil {
		http.Error(w, "Task id not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
