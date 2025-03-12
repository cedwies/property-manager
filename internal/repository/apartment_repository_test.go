package repository

import (
	"database/sql"
	"os"
	"testing"

	"property-management/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDBWithApartments(t *testing.T) (*sql.DB, func()) {
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
	const housesSchema = `
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

	_, err = db.Exec(housesSchema)
	if err != nil {
		t.Fatalf("Failed to create houses schema: %v", err)
	}

	const apartmentsSchema = `
	CREATE TABLE IF NOT EXISTS apartments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		house_id INTEGER NOT NULL,
		size REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (house_id) REFERENCES houses(id)
	);`

	_, err = db.Exec(apartmentsSchema)
	if err != nil {
		t.Fatalf("Failed to create apartments schema: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		db.Close()
		os.Remove(tmpfile.Name())
	}

	return db, cleanup
}

func createTestHouse(t *testing.T, db *sql.DB) *models.House {
	houseRepo := NewHouseRepository(db)
	house := models.NewHouse(
		"Test House",
		"Test Street",
		"123",
		"Test Country",
		"12345",
		"Test City",
	)

	err := houseRepo.Create(house)
	if err != nil {
		t.Fatalf("Failed to create test house: %v", err)
	}

	return house
}

func TestApartmentRepository_Create(t *testing.T) {
	db, cleanup := setupTestDBWithApartments(t)
	defer cleanup()

	house := createTestHouse(t, db)
	repo := NewApartmentRepository(db)

	// Test valid apartment
	apartment, err := models.NewApartment("Test Apartment", house.ID, "75.5")
	if err != nil {
		t.Fatalf("Failed to create apartment model: %v", err)
	}

	err = repo.Create(apartment)
	if err != nil {
		t.Errorf("Error creating apartment: %v", err)
	}

	if apartment.ID == 0 {
		t.Error("Apartment ID should not be 0 after creation")
	}

	// Test with comma as decimal separator
	commaApartment, err := models.NewApartment("Comma Apartment", house.ID, "85,5")
	if err != nil {
		t.Fatalf("Failed to create apartment with comma in size: %v", err)
	}

	err = repo.Create(commaApartment)
	if err != nil {
		t.Errorf("Error creating apartment with comma in size: %v", err)
	}

	// Test invalid apartment (missing name)
	invalidApartment, err := models.NewApartment("", house.ID, "75.5")
	if err == nil {
		err = repo.Create(invalidApartment)
		if err == nil {
			t.Error("Expected error for invalid apartment name, got nil")
		}
	}

	// Test negative size
	_, err = models.NewApartment("Negative Size", house.ID, "-10")
	if err == nil {
		t.Error("Expected error for negative size, got nil")
	}

	// Test zero size
	_, err = models.NewApartment("Zero Size", house.ID, "0")
	if err == nil {
		t.Error("Expected error for zero size, got nil")
	}

	// Test invalid size format
	_, err = models.NewApartment("Invalid Size", house.ID, "abc")
	if err == nil {
		t.Error("Expected error for invalid size format, got nil")
	}

	// Test invalid house ID
	invalidHouseApartment, err := models.NewApartment("Test Apartment", 0, "75.5")
	if err == nil {
		err = repo.Create(invalidHouseApartment)
		if err == nil {
			t.Error("Expected error for invalid house ID, got nil")
		}
	}
}

func TestApartmentRepository_GetAll(t *testing.T) {
	db, cleanup := setupTestDBWithApartments(t)
	defer cleanup()

	house := createTestHouse(t, db)
	repo := NewApartmentRepository(db)

	// Create test apartments
	apartments := []*models.Apartment{}

	apt1, _ := models.NewApartment("Apartment 1", house.ID, "50")
	err := repo.Create(apt1)
	if err != nil {
		t.Fatalf("Error creating test apartment 1: %v", err)
	}
	apartments = append(apartments, apt1)

	apt2, _ := models.NewApartment("Apartment 2", house.ID, "75.5")
	err = repo.Create(apt2)
	if err != nil {
		t.Fatalf("Error creating test apartment 2: %v", err)
	}
	apartments = append(apartments, apt2)

	apt3, _ := models.NewApartment("Apartment 3", house.ID, "100,25")
	err = repo.Create(apt3)
	if err != nil {
		t.Fatalf("Error creating test apartment 3: %v", err)
	}
	apartments = append(apartments, apt3)

	// Test GetAll
	retrievedApartments, err := repo.GetAll()
	if err != nil {
		t.Errorf("Error getting all apartments: %v", err)
	}

	if len(retrievedApartments) != len(apartments) {
		t.Errorf("Expected %d apartments, got %d", len(apartments), len(retrievedApartments))
	}

	// Check that each apartment has its house information
	for _, apartment := range retrievedApartments {
		if apartment.House == nil {
			t.Error("Expected apartment to include house information")
		} else if apartment.House.ID != house.ID {
			t.Errorf("Expected house ID %d, got %d", house.ID, apartment.House.ID)
		}
	}
}

func TestApartmentRepository_GetByHouseID(t *testing.T) {
	db, cleanup := setupTestDBWithApartments(t)
	defer cleanup()

	house1 := createTestHouse(t, db)

	// Create a second house
	houseRepo := NewHouseRepository(db)
	house2 := models.NewHouse(
		"Test House 2",
		"Another Street",
		"456",
		"Another Country",
		"67890",
		"Another City",
	)
	err := houseRepo.Create(house2)
	if err != nil {
		t.Fatalf("Failed to create second test house: %v", err)
	}

	repo := NewApartmentRepository(db)

	// Create apartments for both houses
	apartment1, _ := models.NewApartment("Apartment 1", house1.ID, "50")
	err = repo.Create(apartment1)
	if err != nil {
		t.Fatalf("Error creating test apartment 1: %v", err)
	}

	apartment2, _ := models.NewApartment("Apartment 2", house1.ID, "75.5")
	err = repo.Create(apartment2)
	if err != nil {
		t.Fatalf("Error creating test apartment 2: %v", err)
	}

	apartment3, _ := models.NewApartment("Apartment 3", house2.ID, "100,25")
	err = repo.Create(apartment3)
	if err != nil {
		t.Fatalf("Error creating test apartment 3: %v", err)
	}

	// Test GetByHouseID for house1
	house1Apartments, err := repo.GetByHouseID(house1.ID)
	if err != nil {
		t.Errorf("Error getting apartments for house 1: %v", err)
	}

	if len(house1Apartments) != 2 {
		t.Errorf("Expected 2 apartments for house 1, got %d", len(house1Apartments))
	}

	// Test GetByHouseID for house2
	house2Apartments, err := repo.GetByHouseID(house2.ID)
	if err != nil {
		t.Errorf("Error getting apartments for house 2: %v", err)
	}

	if len(house2Apartments) != 1 {
		t.Errorf("Expected 1 apartment for house 2, got %d", len(house2Apartments))
	}

	// Test GetByHouseID for non-existent house
	nonExistentHouseApartments, err := repo.GetByHouseID(9999)
	if err != nil {
		t.Errorf("Error getting apartments for non-existent house: %v", err)
	}

	if len(nonExistentHouseApartments) != 0 {
		t.Errorf("Expected 0 apartments for non-existent house, got %d", len(nonExistentHouseApartments))
	}
}

func TestApartmentRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDBWithApartments(t)
	defer cleanup()

	house := createTestHouse(t, db)
	repo := NewApartmentRepository(db)

	// Create a test apartment
	apartment, _ := models.NewApartment("Test Apartment", house.ID, "75.5")
	err := repo.Create(apartment)
	if err != nil {
		t.Fatalf("Error creating test apartment: %v", err)
	}

	// Test GetByID with valid ID
	retrievedApartment, err := repo.GetByID(apartment.ID)
	if err != nil {
		t.Errorf("Error getting apartment by ID: %v", err)
	}

	if retrievedApartment.ID != apartment.ID || retrievedApartment.Name != apartment.Name {
		t.Errorf("Retrieved apartment does not match the original")
	}

	if retrievedApartment.House == nil || retrievedApartment.House.ID != house.ID {
		t.Errorf("Retrieved apartment's house information is incorrect")
	}

	// Test GetByID with invalid ID
	_, err = repo.GetByID(9999)
	if err == nil {
		t.Error("Expected error for non-existent apartment, got nil")
	}
}

func TestApartmentRepository_Update(t *testing.T) {
	db, cleanup := setupTestDBWithApartments(t)
	defer cleanup()

	house1 := createTestHouse(t, db)

	// Create a second house
	houseRepo := NewHouseRepository(db)
	house2 := models.NewHouse(
		"Test House 2",
		"Another Street",
		"456",
		"Another Country",
		"67890",
		"Another City",
	)
	err := houseRepo.Create(house2)
	if err != nil {
		t.Fatalf("Failed to create second test house: %v", err)
	}

	repo := NewApartmentRepository(db)

	// Create a test apartment
	apartment, _ := models.NewApartment("Test Apartment", house1.ID, "75.5")
	err = repo.Create(apartment)
	if err != nil {
		t.Fatalf("Error creating test apartment: %v", err)
	}

	// Update the apartment
	apartment.Name = "Updated Apartment"
	apartment.HouseID = house2.ID
	apartment.Size = 100.25

	err = repo.Update(apartment)
	if err != nil {
		t.Errorf("Error updating apartment: %v", err)
	}

	// Verify the update
	retrievedApartment, err := repo.GetByID(apartment.ID)
	if err != nil {
		t.Errorf("Error getting updated apartment: %v", err)
	}

	if retrievedApartment.Name != "Updated Apartment" || retrievedApartment.Size != 100.25 || retrievedApartment.HouseID != house2.ID {
		t.Errorf("Apartment was not properly updated")
	}

	// Test updating with invalid data
	apartment.Name = ""
	err = repo.Update(apartment)
	if err == nil {
		t.Error("Expected error for invalid apartment name, got nil")
	}

	// Test updating non-existent apartment
	nonExistentApartment := &models.Apartment{
		ID:      9999,
		Name:    "Non-existent Apartment",
		HouseID: house1.ID,
		Size:    75.5,
	}

	err = repo.Update(nonExistentApartment)
	if err == nil {
		t.Error("Expected error for updating non-existent apartment, got nil")
	}
}

func TestApartmentRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDBWithApartments(t)
	defer cleanup()

	house := createTestHouse(t, db)
	repo := NewApartmentRepository(db)

	// Create a test apartment
	apartment, _ := models.NewApartment("Test Apartment", house.ID, "75.5")
	err := repo.Create(apartment)
	if err != nil {
		t.Fatalf("Error creating test apartment: %v", err)
	}

	// Delete the apartment
	err = repo.Delete(apartment.ID)
	if err != nil {
		t.Errorf("Error deleting apartment: %v", err)
	}

	// Verify the deletion
	_, err = repo.GetByID(apartment.ID)
	if err == nil {
		t.Error("Expected error when getting deleted apartment, got nil")
	}

	// Test deleting non-existent apartment
	err = repo.Delete(9999)
	if err == nil {
		t.Error("Expected error for deleting non-existent apartment, got nil")
	}
}
