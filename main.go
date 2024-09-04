package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

// Обработчик для PATCH-запроса
func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Логируем начало обработки запроса
	log.Println("PATCH request received")

	// Получаем ID из URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	log.Println("ID from URL:", idStr)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Println("Invalid ID:", idStr)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	log.Println("Parsed ID:", id)

	// Находим сообщение по ID
	var message Message
	if err := DB.First(&message, uint(id)).Error; err != nil {
		log.Println("Message not found with ID:", id)
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	log.Println("Message found:", message)

	// Парсим JSON из тела запроса
	var reqBody requestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Println("Failed to decode request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Request body parsed successfully:", reqBody)

	// Обновляем текст сообщения
	message.Text = reqBody.Message
	log.Println("Updated message text to:", message.Text)

	// Сохраняем обновление в базе данных
	if err := DB.Save(&message).Error; err != nil {
		log.Println("Failed to update message in database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Message updated successfully in the database:", message)

	// Возвращаем обновленное сообщение в ответе
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Println("Failed to encode response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Println("Response successfully sent to client")
	}
}

// Обработчик для DELETE-запроса
func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Логируем начало обработки запроса
	log.Println("DELETE request received")

	// Получаем ID из URL
	vars := mux.Vars(r)
	idStr := vars["id"]
	log.Println("ID from URL:", idStr)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Println("Invalid ID:", idStr)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	log.Println("Parsed ID:", id)

	// Находим сообщение по ID
	var message Message
	if err := DB.First(&message, uint(id)).Error; err != nil {
		log.Println("Message not found with ID:", id)
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	log.Println("Message found:", message)

	// Удаляем сообщение из базы данных
	if err := DB.Delete(&message).Error; err != nil {
		log.Println("Failed to delete message:", err)
		http.Error(w, "Failed to delete message", http.StatusInternalServerError)
		return
	}

	log.Println("Message deleted successfully")

	// Возвращаем подтверждение об удалении
	w.WriteHeader(http.StatusNoContent) // Отправляем 204 No Content
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
	// Маршрут для PATCH-запроса
	router.HandleFunc("/api/message/{id}", UpdateMessageHandler).Methods("PATCH")

	// Маршрут для DELETE-запроса
	router.HandleFunc("/api/message/{id}", DeleteMessageHandler).Methods("DELETE")

	// Запускаем сервер
	http.ListenAndServe(":8080", router)
}
