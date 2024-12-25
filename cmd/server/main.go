package main

import (
	"log"
	"net/http"
	"todolist-api/internal/db"
	"todolist-api/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()
	defer db.CloseDB()

	r := mux.NewRouter()
	r.HandleFunc("/api/tasks", handlers.GetTasks).Methods("GET")
	r.HandleFunc("/api/tasks", handlers.CreateTask).Methods("POST")

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}