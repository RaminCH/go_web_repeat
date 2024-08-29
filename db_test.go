package main

import (
	"testing"
)

// TestInitDB checks if the InitDB function successfully initializes the DB connection.
func TestInitDB(t *testing.T) {
	// Initialize the database connection
	InitDB()

	// Check if the DB variable is initialized
	if DB == nil {
		t.Fatal("Expected DB to be initialized, but it is nil")
	}

	// Test the connection by running a simple query
	var result int
	err := DB.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		t.Fatalf("Failed to execute test query: %v", err)
	}

	// Verify the query result
	if result != 1 {
		t.Fatalf("Expected 1, got %d", result)
	}

	t.Log("Database connection and query test passed successfully")
}
