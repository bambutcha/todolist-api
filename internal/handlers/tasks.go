package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"todolist-api/internal/db"
	"todolist-api/internal/models"

	"github.com/gorilla/mux"
)

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskIDSTR := vars["id"]
	taskID, _ := strconv.Atoi(taskIDSTR)

	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
	}

	query := "UPDATE tasks SET title = ?, description = ?, completed = ? WHERE id = ?"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		http.Error(w, "Failed to prepare query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.Title, task.Description, task.Completed, taskID)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	query := "SELECT id, title, description, completed, created_at FROM tasks WHERE id = ?"
	row := db.DB.QueryRow(query, taskID)

	mskLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		http.Error(w, "Failed to load Moscow time zone: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var task models.Task
	var createdAt []byte
	err = row.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		} else {
			log.Println("Error scanning task:", err)
			http.Error(w, "Failed to get task: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if len(createdAt) > 0 {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err == nil {
			mskTime := parsedTime.In(mskLocation)
			task.CreatedAt = &mskTime
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, title, description, completed, created_at FROM tasks"
	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	mskLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		http.Error(w, "Failed to load Moscow time zone: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var createdAt []byte // Сначала сканируем как срез байтов
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &createdAt)
		if err != nil {
			log.Println("Error scanning task:", err)
			http.Error(w, "Failed to scan task: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Если createdAt не пустое, преобразуем в *time.Time
		if len(createdAt) > 0 {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
			if err == nil {
				mskTime := parsedTime.In(mskLocation)
				task.CreatedAt = &mskTime
			}
		}

		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	stmt, err := db.DB.Prepare("INSERT INTO tasks (title, description) VALUES (?, ?)")
	if err != nil {
		http.Error(w, "Failed to prepare query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(task.Title, task.Description)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get task id", http.StatusInternalServerError)
		return
	}
	task.ID = int(taskID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskIDSTR := vars["id"]
	taskID, _ := strconv.Atoi(taskIDSTR)

	query := "DELETE FROM tasks WHERE id = ?"
	_, err := db.DB.Exec(query, taskID)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
