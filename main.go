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
	var messages []Message

	// Извлекаем все записи из базы данных
	if err := DB.Find(&messages).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Преобразуем слайс сообщений в JSON и отправляем клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
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

	// Создаём новую запись Message с текстом из запроса
	newMessage := Message{Text: reqBody.Message}

	// Сохраняем новое сообщение в базу данных
	if err := DB.Create(&newMessage).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Обновляем глобальную переменную message (если необходимо)
	message = reqBody.Message

	// Отправляем подтверждение клиенту
	fmt.Fprintln(w, "Message received and saved to the database!")
}

func main() {
	// Вызываем метод InitDB() из файла db.go
	InitDB()
	// Автоматическая миграция модели Message
	DB.AutoMigrate(&Message{})

	// Создаём новый маршрутизатор
	router := mux.NewRouter()
	// Маршрут для GET-запроса
	router.HandleFunc("/api/hello", HelloHandler).Methods("GET")
	// Маршрут для POST-запроса
	router.HandleFunc("/api/message", MessageHandler).Methods("POST")
	// Запускаем сервер
	http.ListenAndServe(":8080", router)
}
