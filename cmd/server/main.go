package main

import (
	"log"
	"net/http"
	"todolist-api/internal/db"
	"todolist-api/internal/handlers"
	"github.com/gorilla/mux"
)

func handleFunc() {
	rtp := mux.NewRouter()
	rtp.HandleFunc("/api/tasks", handlers.GetTasks).Methods("GET")
	rtp.HandleFunc("/api/tasks/{id:[0-9]+}", handlers.GetTaskByID).Methods("GET")
	rtp.HandleFunc("/api/tasks", handlers.CreateTask).Methods("POST")
	rtp.HandleFunc("/api/tasks/{id:[0-9]+}", handlers.UpdateTask).Methods("PUT")
	rtp.HandleFunc("/api/tasks/{id:[0-9]+}", handlers.DeleteTask).Methods("DELETE")

	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", rtp)
}

func main() {
	db.InitDB()
	defer db.CloseDB()

	handleFunc()
}