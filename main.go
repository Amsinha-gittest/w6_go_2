package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` // pending or completed
}

var tasks []Task
var idCounter = 1

// Get all tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// Create a new task
func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	_ = json.NewDecoder(r.Body).Decode(&task)

	task.ID = idCounter
	idCounter++
	task.Status = "pending" // default status is pending
	tasks = append(tasks, task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Update a task by ID
func updateTask(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updatedTask Task
	_ = json.NewDecoder(r.Body).Decode(&updatedTask)

	for i, task := range tasks {
		if task.ID == id {
			updatedTask.ID = id // ensure the ID remains the same
			tasks[i] = updatedTask
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedTask)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

// Delete a task by ID
func deleteTask(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent) // 204 No Content
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	// Sample initial task
	tasks = append(tasks, Task{ID: 1, Title: "Sample Task", Description: "This is a sample task", Status: "pending"})

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getTasks(w, r)
		case "POST":
			createTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			updateTask(w, r)
		case "DELETE":
			deleteTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
