package handlers

import (
    "encoding/json"
    "net/http"
    "todolist-api/internal/db"
    "todolist-api/internal/models"
    "log"
    "time"
)

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
