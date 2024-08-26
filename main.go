package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Глобальная переменная для хранения сообщения
var message string

// Структура для парсинга JSON из тела POST-запроса
type requestBody struct {
	Message string `json:"message"`
}

// Обработчик для GET-запроса
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!\n", message)
}

// Обработчик для POST-запроса
func MessageHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody requestBody
	// Парсим JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Обновляем глобальную переменную message
	message = reqBody.Message
	fmt.Fprintln(w, "Message received!")
}

func main() {
	// Создаём новый маршрутизатор
	router := mux.NewRouter()
	// Маршрут для GET-запроса
	router.HandleFunc("/api/hello", HelloHandler).Methods("GET")
	// Маршрут для POST-запроса
	router.HandleFunc("/api/message", MessageHandler).Methods("POST")
	// Запускаем сервер
	http.ListenAndServe(":8080", router)
}
