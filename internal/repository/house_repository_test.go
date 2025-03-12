package repository

import (
	"database/sql"
	"os"
	"testing"

	"property-management/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	// Create a temporary file for the test database
	tmpfile, err := os.CreateTemp("", "test_property_management.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpfile.Close()

	// Open the database connection
	db, err := sql.Open("sqlite3", tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	// Create the schema
	const schema = `
	CREATE TABLE IF NOT EXISTS houses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		street TEXT NOT NULL,
		number TEXT NOT NULL,
		country TEXT NOT NULL,
		zip_code TEXT NOT NULL,
		city TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("Failed to create schema: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		db.Close()
		os.Remove(tmpfile.Name())
	}

	return db, cleanup
}

func TestHouseRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewHouseRepository(db)

	// Test valid house
	house := models.NewHouse(
		"Test House",
		"Test Street",
		"123",
		"Test Country",
		"12345",
		"Test City",
	)

	err := repo.Create(house)
	if err != nil {
		t.Errorf("Error creating house: %v", err)
	}

	if house.ID == 0 {
		t.Error("House ID should not be 0 after creation")
	}

	// Test invalid house (missing name)
	invalidHouse := models.NewHouse(
		"", // Empty name
		"Test Street",
		"123",
		"Test Country",
		"12345",
		"Test City",
	)

	err = repo.Create(invalidHouse)
	if err == nil {
		t.Error("Expected error for invalid house, got nil")
	}

	// Test invalid city (just a number)
	invalidCityHouse := models.NewHouse(
		"Test House",
		"Test Street",
		"123",
		"Test Country",
		"12345",
		"123", // City is just a number
	)

	err = repo.Create(invalidCityHouse)
	if err == nil {
		t.Error("Expected error for invalid city, got nil")
	}
}

func TestHouseRepository_GetAll(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewHouseRepository(db)

	// Create test houses
	houses := []*models.House{
		models.NewHouse("House 1", "Street 1", "1", "Country 1", "11111", "City 1"),
		models.NewHouse("House 2", "Street 2", "2", "Country 2", "22222", "City 2"),
		models.NewHouse("House 3", "Street 3", "3", "Country 3", "33333", "City 3"),
	}

	for _, house := range houses {
		err := repo.Create(house)
		if err != nil {
			t.Fatalf("Error creating test house: %v", err)
		}
	}

	// Test GetAll
	retrievedHouses, err := repo.GetAll()
	if err != nil {
		t.Errorf("Error getting all houses: %v", err)
	}

	if len(retrievedHouses) != len(houses) {
		t.Errorf("Expected %d houses, got %d", len(houses), len(retrievedHouses))
	}
}

func TestHouseRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewHouseRepository(db)

	// Create a test house
	house := models.NewHouse(
		"Test House",
		"Test Street",
		"123",
		"Test Country",
		"12345",
		"Test City",
	)

	err := repo.Create(house)
	if err != nil {
		t.Fatalf("Error creating test house: %v", err)
	}

	// Test GetByID with valid ID
	retrievedHouse, err := repo.GetByID(house.ID)
	if err != nil {
		t.Errorf("Error getting house by ID: %v", err)
	}

	if retrievedHouse.ID != house.ID || retrievedHouse.Name != house.Name {
		t.Errorf("Retrieved house does not match the original")
	}

	// Test GetByID with invalid ID
	_, err = repo.GetByID(9999)
	if err == nil {
		t.Error("Expected error for non-existent house, got nil")
	}
}

func TestHouseRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewHouseRepository(db)

	// Create a test house
	house := models.NewHouse(
		"Test House",
		"Test Street",
		"123",
		"Test Country",
		"12345",
		"Test City",
	)

	err := repo.Create(house)
	if err != nil {
		t.Fatalf("Error creating test house: %v", err)
	}

	// Update the house
	house.Name = "Updated House"
	house.Street = "Updated Street"

	err = repo.Update(house)
	if err != nil {
		t.Errorf("Error updating house: %v", err)
	}

	// Verify the update
	retrievedHouse, err := repo.GetByID(house.ID)
	if err != nil {
		t.Errorf("Error getting updated house: %v", err)
	}

	if retrievedHouse.Name != "Updated House" || retrievedHouse.Street != "Updated Street" {
		t.Errorf("House was not properly updated")
	}

	// Test updating non-existent house
	nonExistentHouse := models.NewHouse(
		"Non-existent House",
		"Street",
		"123",
		"Country",
		"12345",
		"City",
	)
	nonExistentHouse.ID = 9999

	err = repo.Update(nonExistentHouse)
	if err == nil {
		t.Error("Expected error for updating non-existent house, got nil")
	}
}

func TestHouseRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewHouseRepository(db)

	// Create a test house
	house := models.NewHouse(
		"Test House",
		"Test Street",
		"123",
		"Test Country",
		"12345",
		"Test City",
	)

	err := repo.Create(house)
	if err != nil {
		t.Fatalf("Error creating test house: %v", err)
	}

	// Delete the house
	err = repo.Delete(house.ID)
	if err != nil {
		t.Errorf("Error deleting house: %v", err)
	}

	// Verify the deletion
	_, err = repo.GetByID(house.ID)
	if err == nil {
		t.Error("Expected error when getting deleted house, got nil")
	}

	// Test deleting non-existent house
	err = repo.Delete(9999)
	if err == nil {
		t.Error("Expected error for deleting non-existent house, got nil")
	}
}
