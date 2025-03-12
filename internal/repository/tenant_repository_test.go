package repository

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"property-management/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDBWithTenants(t *testing.T) (*sql.DB, func()) {
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

	const tenantsSchema = `
	CREATE TABLE IF NOT EXISTS tenants (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		move_in_date TEXT NOT NULL,
		move_out_date TEXT,
		deposit REAL NOT NULL,
		email TEXT,
		number_of_persons INTEGER NOT NULL,
		target_cold_rent REAL NOT NULL,
		target_ancillary_payment REAL NOT NULL,
		target_electricity_payment REAL NOT NULL,
		greeting TEXT NOT NULL,
		house_id INTEGER NOT NULL,
		apartment_id INTEGER NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (house_id) REFERENCES houses(id),
		FOREIGN KEY (apartment_id) REFERENCES apartments(id)
	);`

	_, err = db.Exec(tenantsSchema)
	if err != nil {
		t.Fatalf("Failed to create tenants schema: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		db.Close()
		os.Remove(tmpfile.Name())
	}

	return db, cleanup
}

func createTestHouseAndApartment(t *testing.T, db *sql.DB) (*models.House, *models.Apartment) {
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

	// Create a test apartment
	apartmentRepo := NewApartmentRepository(db)
	apartment, err := models.NewApartment("Test Apartment", house.ID, "75.5")
	if err != nil {
		t.Fatalf("Failed to create apartment model: %v", err)
	}

	err = apartmentRepo.Create(apartment)
	if err != nil {
		t.Fatalf("Failed to create test apartment: %v", err)
	}

	return house, apartment
}

func createSecondHouseAndApartment(t *testing.T, db *sql.DB) (*models.House, *models.Apartment) {
	houseRepo := NewHouseRepository(db)
	house := models.NewHouse(
		"Second House",
		"Second Street",
		"456",
		"Second Country",
		"67890",
		"Second City",
	)

	err := houseRepo.Create(house)
	if err != nil {
		t.Fatalf("Failed to create second test house: %v", err)
	}

	// Create a test apartment
	apartmentRepo := NewApartmentRepository(db)
	apartment, err := models.NewApartment("Second Apartment", house.ID, "100")
	if err != nil {
		t.Fatalf("Failed to create second apartment model: %v", err)
	}

	err = apartmentRepo.Create(apartment)
	if err != nil {
		t.Fatalf("Failed to create second test apartment: %v", err)
	}

	return house, apartment
}

func createValidTenant(t *testing.T, house *models.House, apartment *models.Apartment) *models.Tenant {
	tenant, err := models.NewTenant(
		"John",                 // FirstName
		"Doe",                  // LastName
		"2023-01-01",           // MoveInDate
		"",                     // MoveOutDate (empty for current tenant)
		"1500",                 // Deposit
		"john.doe@example.com", // Email
		"2",                    // NumberOfPersons
		"800",                  // TargetColdRent
		"200",                  // TargetAncillaryPayment
		"100",                  // TargetElectricityPayment
		"Dear Mr. Doe",         // Greeting
		house.ID,               // HouseID
		apartment.ID,           // ApartmentID
	)
	if err != nil {
		t.Fatalf("Failed to create valid tenant: %v", err)
	}
	return tenant
}

func TestTenantRepository_Create(t *testing.T) {
	db, cleanup := setupTestDBWithTenants(t)
	defer cleanup()

	house, apartment := createTestHouseAndApartment(t, db)
	repo := NewTenantRepository(db)

	// Test valid tenant
	tenant := createValidTenant(t, house, apartment)

	err := repo.Create(tenant)
	if err != nil {
		t.Errorf("Error creating tenant: %v", err)
	}

	if tenant.ID == 0 {
		t.Error("Tenant ID should not be 0 after creation")
	}

	// Test with move-out date
	tenantWithMoveOutDate, err := models.NewTenant(
		"Jane",
		"Smith",
		"2022-01-01",
		"2023-01-01",
		"1200",
		"jane.smith@example.com",
		"1",
		"600",
		"150",
		"80",
		"Dear Ms. Smith",
		house.ID,
		apartment.ID,
	)
	if err != nil {
		t.Fatalf("Failed to create tenant with move-out date: %v", err)
	}

	err = repo.Create(tenantWithMoveOutDate)
	if err != nil {
		t.Errorf("Error creating tenant with move-out date: %v", err)
	}

	// Test invalid tenant (move-out date before move-in date)
	_, err = models.NewTenant(
		"Invalid",
		"Tenant",
		"2023-01-01",
		"2022-01-01", // Earlier than move-in date
		"1000",
		"invalid@example.com",
		"1",
		"500",
		"100",
		"50",
		"Dear Tenant",
		house.ID,
		apartment.ID,
	)
	if err == nil {
		t.Error("Expected error for move-out date before move-in date, got nil")
	}

	// Test invalid number of persons
	_, err = models.NewTenant(
		"Invalid",
		"Tenant",
		"2023-01-01",
		"",
		"1000",
		"invalid@example.com",
		"0", // Invalid number of persons
		"500",
		"100",
		"50",
		"Dear Tenant",
		house.ID,
		apartment.ID,
	)
	if err == nil {
		t.Error("Expected error for invalid number of persons, got nil")
	}

	// Test invalid target cold rent
	_, err = models.NewTenant(
		"Invalid",
		"Tenant",
		"2023-01-01",
		"",
		"1000",
		"invalid@example.com",
		"1",
		"0", // Invalid cold rent
		"100",
		"50",
		"Dear Tenant",
		house.ID,
		apartment.ID,
	)
	if err == nil {
		t.Error("Expected error for invalid target cold rent, got nil")
	}

	// Test invalid apartment (not belonging to the house)
	_, apartment2 := createSecondHouseAndApartment(t, db)

	invalidTenant, err := models.NewTenant(
		"Invalid",
		"Tenant",
		"2023-01-01",
		"",
		"1000",
		"invalid@example.com",
		"1",
		"500",
		"100",
		"50",
		"Dear Tenant",
		house.ID,      // First house
		apartment2.ID, // Second apartment (belongs to second house)
	)
	if err == nil {
		err = repo.Create(invalidTenant)
		if err == nil {
			t.Error("Expected error for apartment not belonging to house, got nil")
		}
	}

	// Test invalid house ID
	_, err = models.NewTenant(
		"Invalid",
		"Tenant",
		"2023-01-01",
		"",
		"1000",
		"invalid@example.com",
		"1",
		"500",
		"100",
		"50",
		"Dear Tenant",
		0, // Invalid house ID
		apartment.ID,
	)
	if err == nil {
		t.Error("Expected error for invalid house ID, got nil")
	}

	// Test invalid apartment ID
	_, err = models.NewTenant(
		"Invalid",
		"Tenant",
		"2023-01-01",
		"",
		"1000",
		"invalid@example.com",
		"1",
		"500",
		"100",
		"50",
		"Dear Tenant",
		house.ID,
		0, // Invalid apartment ID
	)
	if err == nil {
		t.Error("Expected error for invalid apartment ID, got nil")
	}
}

func TestTenantRepository_GetAll(t *testing.T) {
	db, cleanup := setupTestDBWithTenants(t)
	defer cleanup()

	house, apartment := createTestHouseAndApartment(t, db)
	house2, apartment2 := createSecondHouseAndApartment(t, db)
	repo := NewTenantRepository(db)

	// Create test tenants
	tenants := []*models.Tenant{}

	tenant1 := createValidTenant(t, house, apartment)
	err := repo.Create(tenant1)
	if err != nil {
		t.Fatalf("Error creating test tenant 1: %v", err)
	}
	tenants = append(tenants, tenant1)

	tenant2, err := models.NewTenant(
		"Jane",
		"Smith",
		"2022-01-01",
		"2023-01-01",
		"1200",
		"jane.smith@example.com",
		"1",
		"600",
		"150",
		"80",
		"Dear Ms. Smith",
		house.ID,
		apartment.ID,
	)
	if err != nil {
		t.Fatalf("Error creating tenant model 2: %v", err)
	}
	err = repo.Create(tenant2)
	if err != nil {
		t.Fatalf("Error creating test tenant 2: %v", err)
	}
	tenants = append(tenants, tenant2)

	tenant3, err := models.NewTenant(
		"Alice",
		"Johnson",
		"2023-02-01",
		"",
		"2000",
		"alice.johnson@example.com",
		"3",
		"900",
		"250",
		"120",
		"Dear Ms. Johnson",
		house2.ID,
		apartment2.ID,
	)
	if err != nil {
		t.Fatalf("Error creating tenant model 3: %v", err)
	}
	err = repo.Create(tenant3)
	if err != nil {
		t.Fatalf("Error creating test tenant 3: %v", err)
	}
	tenants = append(tenants, tenant3)

	// Test GetAll
	retrievedTenants, err := repo.GetAll()
	if err != nil {
		t.Errorf("Error getting all tenants: %v", err)
	}

	if len(retrievedTenants) != len(tenants) {
		t.Errorf("Expected %d tenants, got %d", len(tenants), len(retrievedTenants))
	}

	// Check that each tenant has house and apartment information
	for _, tenant := range retrievedTenants {
		if tenant.House == nil {
			t.Error("Expected tenant to include house information")
		}
		if tenant.Apartment == nil {
			t.Error("Expected tenant to include apartment information")
		}
	}
}

func TestTenantRepository_GetByHouseID(t *testing.T) {
	db, cleanup := setupTestDBWithTenants(t)
	defer cleanup()

	house1, apartment1 := createTestHouseAndApartment(t, db)
	house2, apartment2 := createSecondHouseAndApartment(t, db)
	repo := NewTenantRepository(db)

	// Create tenants for both houses
	tenant1 := createValidTenant(t, house1, apartment1)
	err := repo.Create(tenant1)
	if err != nil {
		t.Fatalf("Error creating test tenant 1: %v", err)
	}

	tenant2, err := models.NewTenant(
		"Jane",
		"Smith",
		"2022-01-01",
		"2023-01-01",
		"1200",
		"jane.smith@example.com",
		"1",
		"600",
		"150",
		"80",
		"Dear Ms. Smith",
		house1.ID,
		apartment1.ID,
	)
	if err != nil {
		t.Fatalf("Error creating tenant model 2: %v", err)
	}
	err = repo.Create(tenant2)
	if err != nil {
		t.Fatalf("Error creating test tenant 2: %v", err)
	}

	tenant3, err := models.NewTenant(
		"Alice",
		"Johnson",
		"2023-02-01",
		"",
		"2000",
		"alice.johnson@example.com",
		"3",
		"900",
		"250",
		"120",
		"Dear Ms. Johnson",
		house2.ID,
		apartment2.ID,
	)
	if err != nil {
		t.Fatalf("Error creating tenant model 3: %v", err)
	}
	err = repo.Create(tenant3)
	if err != nil {
		t.Fatalf("Error creating test tenant 3: %v", err)
	}

	// Test GetByHouseID for house1
	house1Tenants, err := repo.GetByHouseID(house1.ID)
	if err != nil {
		t.Errorf("Error getting tenants for house 1: %v", err)
	}

	if len(house1Tenants) != 2 {
		t.Errorf("Expected 2 tenants for house 1, got %d", len(house1Tenants))
	}

	// Test GetByHouseID for house2
	house2Tenants, err := repo.GetByHouseID(house2.ID)
	if err != nil {
		t.Errorf("Error getting tenants for house 2: %v", err)
	}

	if len(house2Tenants) != 1 {
		t.Errorf("Expected 1 tenant for house 2, got %d", len(house2Tenants))
	}

	// Test GetByHouseID for non-existent house
	nonExistentHouseTenants, err := repo.GetByHouseID(9999)
	if err != nil {
		t.Errorf("Error getting tenants for non-existent house: %v", err)
	}

	if len(nonExistentHouseTenants) != 0 {
		t.Errorf("Expected 0 tenants for non-existent house, got %d", len(nonExistentHouseTenants))
	}
}

func TestTenantRepository_GetByApartmentID(t *testing.T) {
	db, cleanup := setupTestDBWithTenants(t)
	defer cleanup()

	house1, apartment1 := createTestHouseAndApartment(t, db)

	// Create a second apartment in the same house
	apartmentRepo := NewApartmentRepository(db)
	apartment2, err := models.NewApartment("Second Apartment", house1.ID, "100")
	if err != nil {
		t.Fatalf("Failed to create second apartment model: %v", err)
	}
	err = apartmentRepo.Create(apartment2)
	if err != nil {
		t.Fatalf("Failed to create second test apartment: %v", err)
	}

	tenantRepo := NewTenantRepository(db)

	// Create tenants for both apartments
	tenant1 := createValidTenant(t, house1, apartment1)
	err = tenantRepo.Create(tenant1)
	if err != nil {
		t.Fatalf("Error creating test tenant 1: %v", err)
	}

	tenant2, err := models.NewTenant(
		"Jane",
		"Smith",
		"2022-01-01",
		"2023-01-01",
		"1200",
		"jane.smith@example.com",
		"1",
		"600",
		"150",
		"80",
		"Dear Ms. Smith",
		house1.ID,
		apartment2.ID,
	)
	if err != nil {
		t.Fatalf("Error creating tenant model 2: %v", err)
	}
	err = tenantRepo.Create(tenant2)
	if err != nil {
		t.Fatalf("Error creating test tenant 2: %v", err)
	}

	// Test GetByApartmentID for apartment1
	apartment1Tenants, err := tenantRepo.GetByApartmentID(apartment1.ID)
	if err != nil {
		t.Errorf("Error getting tenants for apartment 1: %v", err)
	}

	if len(apartment1Tenants) != 1 {
		t.Errorf("Expected 1 tenant for apartment 1, got %d", len(apartment1Tenants))
	}

	// Test GetByApartmentID for apartment2
	apartment2Tenants, err := tenantRepo.GetByApartmentID(apartment2.ID)
	if err != nil {
		t.Errorf("Error getting tenants for apartment 2: %v", err)
	}

	if len(apartment2Tenants) != 1 {
		t.Errorf("Expected 1 tenant for apartment 2, got %d", len(apartment2Tenants))
	}

	// Test GetByApartmentID for non-existent apartment
	nonExistentApartmentTenants, err := tenantRepo.GetByApartmentID(9999)
	if err != nil {
		t.Errorf("Error getting tenants for non-existent apartment: %v", err)
	}

	if len(nonExistentApartmentTenants) != 0 {
		t.Errorf("Expected 0 tenants for non-existent apartment, got %d", len(nonExistentApartmentTenants))
	}
}

func TestTenantRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDBWithTenants(t)
	defer cleanup()

	house, apartment := createTestHouseAndApartment(t, db)
	repo := NewTenantRepository(db)

	// Create a test tenant
	tenant := createValidTenant(t, house, apartment)
	err := repo.Create(tenant)
	if err != nil {
		t.Fatalf("Error creating test tenant: %v", err)
	}

	// Test GetByID with valid ID
	retrievedTenant, err := repo.GetByID(tenant.ID)
	if err != nil {
		t.Errorf("Error getting tenant by ID: %v", err)
	}

	if retrievedTenant.ID != tenant.ID ||
		retrievedTenant.FirstName != tenant.FirstName ||
		retrievedTenant.LastName != tenant.LastName {
		t.Errorf("Retrieved tenant does not match the original")
	}

	if retrievedTenant.House == nil || retrievedTenant.House.ID != house.ID {
		t.Errorf("Retrieved tenant's house information is incorrect")
	}

	if retrievedTenant.Apartment == nil || retrievedTenant.Apartment.ID != apartment.ID {
		t.Errorf("Retrieved tenant's apartment information is incorrect")
	}

	// Test GetByID with invalid ID
	_, err = repo.GetByID(9999)
	if err == nil {
		t.Error("Expected error for non-existent tenant, got nil")
	}
}

func TestTenantRepository_Update(t *testing.T) {
	db, cleanup := setupTestDBWithTenants(t)
	defer cleanup()

	house1, apartment1 := createTestHouseAndApartment(t, db)
	_, apartment2 := createSecondHouseAndApartment(t, db)
	repo := NewTenantRepository(db)

	// Create a test tenant
	tenant := createValidTenant(t, house1, apartment1)
	err := repo.Create(tenant)
	if err != nil {
		t.Fatalf("Error creating test tenant: %v", err)
	}

	// Update the tenant
	tenant.FirstName = "Updated"
	tenant.LastName = "Name"
	tenant.NumberOfPersons = 3
	tenant.TargetColdRent = 850

	// Valid move to another apartment that belongs to the same house
	apartmentRepo := NewApartmentRepository(db)
	newApartment, err := models.NewApartment("New Apartment", house1.ID, "120")
	if err != nil {
		t.Fatalf("Failed to create new apartment model: %v", err)
	}
	err = apartmentRepo.Create(newApartment)
	if err != nil {
		t.Fatalf("Failed to create new test apartment: %v", err)
	}

	tenant.ApartmentID = newApartment.ID

	err = repo.Update(tenant)
	if err != nil {
		t.Errorf("Error updating tenant: %v", err)
	}

	// Verify the update
	retrievedTenant, err := repo.GetByID(tenant.ID)
	if err != nil {
		t.Errorf("Error getting updated tenant: %v", err)
	}

	if retrievedTenant.FirstName != "Updated" ||
		retrievedTenant.LastName != "Name" ||
		retrievedTenant.NumberOfPersons != 3 ||
		retrievedTenant.TargetColdRent != 850 ||
		retrievedTenant.ApartmentID != newApartment.ID {
		t.Errorf("Tenant was not properly updated")
	}

	// Test updating with invalid data (zero persons)
	tenant.NumberOfPersons = 0
	err = repo.Update(tenant)
	if err == nil {
		t.Error("Expected error for invalid number of persons, got nil")
	}

	// Reset to valid value
	tenant.NumberOfPersons = 1

	// Test updating with apartment from different house
	tenant.HouseID = house1.ID
	tenant.ApartmentID = apartment2.ID // This belongs to house2
	err = repo.Update(tenant)
	if err == nil {
		t.Error("Expected error for apartment not belonging to house, got nil")
	}

	// Test updating non-existent tenant
	nonExistentTenant := models.Tenant{
		ID:                       9999,
		FirstName:                "Non-existent",
		LastName:                 "Tenant",
		MoveInDate:               time.Now(),
		Deposit:                  1000,
		NumberOfPersons:          1,
		TargetColdRent:           500,
		TargetAncillaryPayment:   100,
		TargetElectricityPayment: 50,
		Greeting:                 "Dear Tenant",
		HouseID:                  house1.ID,
		ApartmentID:              apartment1.ID,
	}

	err = repo.Update(&nonExistentTenant)
	if err == nil {
		t.Error("Expected error for updating non-existent tenant, got nil")
	}
}

func TestTenantRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDBWithTenants(t)
	defer cleanup()

	house, apartment := createTestHouseAndApartment(t, db)
	repo := NewTenantRepository(db)

	// Create a test tenant
	tenant := createValidTenant(t, house, apartment)
	err := repo.Create(tenant)
	if err != nil {
		t.Fatalf("Error creating test tenant: %v", err)
	}

	// Delete the tenant
	err = repo.Delete(tenant.ID)
	if err != nil {
		t.Errorf("Error deleting tenant: %v", err)
	}

	// Verify the deletion
	_, err = repo.GetByID(tenant.ID)
	if err == nil {
		t.Error("Expected error when getting deleted tenant, got nil")
	}

	// Test deleting non-existent tenant
	err = repo.Delete(9999)
	if err == nil {
		t.Error("Expected error for deleting non-existent tenant, got nil")
	}
}
