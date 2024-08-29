package main

import (
	"gorm.io/gorm"
)

// Message struct represents a message in the database
type Message struct {
	gorm.Model
	Text string `json:"text"`
}
