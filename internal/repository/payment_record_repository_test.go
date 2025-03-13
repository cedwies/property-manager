package repository

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"property-management/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDBWithPaymentRecords(t *testing.T) (*sql.DB, func()) {
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

	// Create the schema for houses, apartments, tenants, and payment records
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

	const paymentRecordsSchema = `
	CREATE TABLE IF NOT EXISTS payment_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tenant_id INTEGER NOT NULL,
		month TEXT NOT NULL,
		target_cold_rent REAL NOT NULL,
		paid_cold_rent REAL,
		paid_ancillary REAL,
		paid_electricity REAL,
		extra_payments REAL,
		persons INTEGER, 
		note TEXT,
		is_locked BOOLEAN DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (tenant_id) REFERENCES tenants(id),
		UNIQUE(tenant_id, month)
	);`

	_, err = db.Exec(paymentRecordsSchema)
	if err != nil {
		t.Fatalf("Failed to create payment_records schema: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		db.Close()
		os.Remove(tmpfile.Name())
	}

	return db, cleanup
}

func createTestTenant(t *testing.T, db *sql.DB) *models.Tenant {
	house, apartment := createTestHouseAndApartment(t, db)

	// Create a test tenant
	tenantRepo := NewTenantRepository(db)
	tenant, err := models.NewTenant(
		"John",
		"Doe",
		"2023-01-01", // Move-in date
		"",           // No move-out date (current tenant)
		"1500",       // Deposit
		"john.doe@example.com",
		"2",   // Persons
		"800", // Target cold rent
		"200", // Target ancillary
		"100", // Target electricity
		"Dear Mr. Doe",
		house.ID,
		apartment.ID,
	)
	if err != nil {
		t.Fatalf("Failed to create tenant model: %v", err)
	}

	err = tenantRepo.Create(tenant)
	if err != nil {
		t.Fatalf("Failed to create test tenant: %v", err)
	}

	return tenant
}

func createFormerTenant(t *testing.T, db *sql.DB) *models.Tenant {
	house, apartment := createTestHouseAndApartment(t, db)

	// Create a test former tenant
	tenantRepo := NewTenantRepository(db)
	tenant, err := models.NewTenant(
		"Jane",
		"Smith",
		"2022-01-01", // Move-in date
		"2023-01-01", // Move-out date (former tenant)
		"1200",       // Deposit
		"jane.smith@example.com",
		"1",   // Persons
		"600", // Target cold rent
		"150", // Target ancillary
		"80",  // Target electricity
		"Dear Ms. Smith",
		house.ID,
		apartment.ID,
	)
	if err != nil {
		t.Fatalf("Failed to create former tenant model: %v", err)
	}

	err = tenantRepo.Create(tenant)
	if err != nil {
		t.Fatalf("Failed to create test former tenant: %v", err)
	}

	return tenant
}

func TestPaymentRecordRepository_Create(t *testing.T) {
	db, cleanup := setupTestDBWithPaymentRecords(t)
	defer cleanup()

	tenant := createTestTenant(t, db)
	repo := NewPaymentRecordRepository(db)

	// Test valid payment record
	record, err := models.NewPaymentRecord(
		tenant.ID,
		"2023-01", // Month (January 2023)
		tenant.TargetColdRent,
		"800", // Paid cold rent
		"200", // Paid ancillary
		"100", // Paid electricity
		"50",  // Extra payments
		"2",   // Persons
		"Payment received on time",
		false, // Not locked
	)
	if err != nil {
		t.Fatalf("Failed to create payment record model: %v", err)
	}

	err = repo.Create(record)
	if err != nil {
		t.Errorf("Error creating payment record: %v", err)
	}

	if record.ID == 0 {
		t.Error("Payment record ID should not be 0 after creation")
	}

	// Test creating a duplicate record (same tenant and month)
	duplicateRecord, err := models.NewPaymentRecord(
		tenant.ID,
		"2023-01", // Same month
		tenant.TargetColdRent,
		"850",
		"210",
		"110",
		"0",
		"2",
		"Another payment",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to create duplicate payment record model: %v", err)
	}

	err = repo.Create(duplicateRecord)
	if err == nil {
		t.Error("Expected error for duplicate payment record, got nil")
	}

	// Test with invalid month format
	_, err = models.NewPaymentRecord(
		tenant.ID,
		"01-2023", // Invalid format
		tenant.TargetColdRent,
		"800",
		"200",
		"100",
		"50",
		"2",
		"Note",
		false,
	)
	if err == nil {
		t.Error("Expected error for invalid month format, got nil")
	}

	// Test with negative payment amount
	_, err = models.NewPaymentRecord(
		tenant.ID,
		"2023-02",
		tenant.TargetColdRent,
		"-100", // Negative payment
		"200",
		"100",
		"50",
		"2",
		"Note",
		false,
	)
	if err == nil {
		t.Error("Expected error for negative payment amount, got nil")
	}

	// Test with invalid number of persons
	_, err = models.NewPaymentRecord(
		tenant.ID,
		"2023-02",
		tenant.TargetColdRent,
		"800",
		"200",
		"100",
		"50",
		"0", // Invalid persons count
		"Note",
		false,
	)
	if err == nil {
		t.Error("Expected error for invalid number of persons, got nil")
	}

	// Test with non-existent tenant ID
	invalidRecord, err := models.NewPaymentRecord(
		9999, // Non-existent tenant ID
		"2023-02",
		tenant.TargetColdRent,
		"800",
		"200",
		"100",
		"50",
		"2",
		"Note",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to create payment record model with invalid tenant ID: %v", err)
	}

	err = repo.Create(invalidRecord)
	if err == nil {
		t.Error("Expected error for non-existent tenant ID, got nil")
	}
}

func TestPaymentRecordRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDBWithPaymentRecords(t)
	defer cleanup()

	tenant := createTestTenant(t, db)
	repo := NewPaymentRecordRepository(db)

	// Create a test payment record
	record, err := models.NewPaymentRecord(
		tenant.ID,
		"2023-01",
		tenant.TargetColdRent,
		"800",
		"200",
		"100",
		"50",
		"2",
		"Payment received on time",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to create payment record model: %v", err)
	}

	err = repo.Create(record)
	if err != nil {
		t.Fatalf("Error creating test payment record: %v", err)
	}

	// Test GetByID with valid ID
	retrievedRecord, err := repo.GetByID(record.ID)
	if err != nil {
		t.Errorf("Error getting payment record by ID: %v", err)
	}

	if retrievedRecord.ID != record.ID || retrievedRecord.Month != record.Month {
		t.Errorf("Retrieved payment record does not match the original")
	}

	// Test GetByID with invalid ID
	_, err = repo.GetByID(9999)
	if err == nil {
		t.Error("Expected error for non-existent payment record, got nil")
	}
}

func TestPaymentRecordRepository_GetByTenantAndMonth(t *testing.T) {
	db, cleanup := setupTestDBWithPaymentRecords(t)
	defer cleanup()

	tenant := createTestTenant(t, db)
	repo := NewPaymentRecordRepository(db)

	// Create a test payment record
	record, err := models.NewPaymentRecord(
		tenant.ID,
		"2023-01",
		tenant.TargetColdRent,
		"800",
		"200",
		"100",
		"50",
		"2",
		"Payment received on time",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to create payment record model: %v", err)
	}

	err = repo.Create(record)
	if err != nil {
		t.Fatalf("Error creating test payment record: %v", err)
	}

	// Test GetByTenantAndMonth with valid data
	retrievedRecord, err := repo.GetByTenantAndMonth(tenant.ID, "2023-01")
	if err != nil {
		t.Errorf("Error getting payment record by tenant and month: %v", err)
	}

	if retrievedRecord.ID != record.ID || retrievedRecord.Month != record.Month {
		t.Errorf("Retrieved payment record does not match the original")
	}

	// Test with non-existent tenant ID
	_, err = repo.GetByTenantAndMonth(9999, "2023-01")
	if err == nil {
		t.Error("Expected error for non-existent tenant ID, got nil")
	}

	// Test with non-existent month
	_, err = repo.GetByTenantAndMonth(tenant.ID, "2023-02")
	if err == nil {
		t.Error("Expected error for non-existent month, got nil")
	}
}

func TestPaymentRecordRepository_GetByTenantID(t *testing.T) {
	db, cleanup := setupTestDBWithPaymentRecords(t)
	defer cleanup()

	tenant := createTestTenant(t, db)
	repo := NewPaymentRecordRepository(db)

	// Create multiple payment records for the same tenant
	months := []string{"2023-01", "2023-02", "2023-03"}
	for _, month := range months {
		record, err := models.NewPaymentRecord(
			tenant.ID,
			month,
			tenant.TargetColdRent,
			"800",
			"200",
			"100",
			"50",
			"2",
			"Payment for "+month,
			false,
		)
		if err != nil {
			t.Fatalf("Failed to create payment record model for %s: %v", month, err)
		}

		err = repo.Create(record)
		if err != nil {
			t.Fatalf("Error creating test payment record for %s: %v", month, err)
		}
	}

	// Test GetByTenantID
	records, err := repo.GetByTenantID(tenant.ID)
	if err != nil {
		t.Errorf("Error getting payment records by tenant ID: %v", err)
	}

	if len(records) != len(months) {
		t.Errorf("Expected %d payment records, got %d", len(months), len(records))
	}

	// Test with non-existent tenant ID
	records, err = repo.GetByTenantID(9999)
	if err != nil {
		t.Errorf("Unexpected error for non-existent tenant ID: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("Expected 0 payment records for non-existent tenant, got %d", len(records))
	}
}

func TestPaymentRecordRepository_Update(t *testing.T) {
	db, cleanup := setupTestDBWithPaymentRecords(t)
	defer cleanup()

	tenant := createTestTenant(t, db)
	repo := NewPaymentRecordRepository(db)

	// Create a test payment record
	record, err := models.NewPaymentRecord(
		tenant.ID,
		"2023-01",
		tenant.TargetColdRent,
		"800",
		"200",
		"100",
		"50",
		"2",
		"Initial payment",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to create payment record model: %v", err)
	}

	err = repo.Create(record)
	if err != nil {
		t.Fatalf("Error creating test payment record: %v", err)
	}

	// Update the payment record
	record.PaidColdRent = 850
	record.PaidAncillary = 220
	record.PaidElectricity = 110
	record.ExtraPayments = 0
	record.Note = "Updated payment"

	err = repo.Update(record)
	if err != nil {
		t.Errorf("Error updating payment record: %v", err)
	}

	// Verify the update
	retrievedRecord, err := repo.GetByID(record.ID)
	if err != nil {
		t.Errorf("Error getting updated payment record: %v", err)
	}

	if retrievedRecord.PaidColdRent != 850 ||
		retrievedRecord.PaidAncillary != 220 ||
		retrievedRecord.PaidElectricity != 110 ||
		retrievedRecord.ExtraPayments != 0 ||
		retrievedRecord.Note != "Updated payment" {
		t.Errorf("Payment record was not properly updated")
	}

	// Test locking a record
	record.IsLocked = true
	err = repo.Update(record)
	if err != nil {
		t.Errorf("Error locking payment record: %v", err)
	}

	// Verify the record is locked
	retrievedRecord, err = repo.GetByID(record.ID)
	if err != nil {
		t.Errorf("Error getting locked payment record: %v", err)
	}
	if !retrievedRecord.IsLocked {
		t.Errorf("Payment record is not locked")
	}

	// Try to update a locked record
	originalPaidColdRent := retrievedRecord.PaidColdRent
	record.PaidColdRent = 900
	err = repo.Update(record)
	if err != nil {
		t.Errorf("Error trying to update locked payment record: %v", err)
	}

	// Verify the payment values didn't change (except note)
	retrievedRecord, err = repo.GetByID(record.ID)
	if err != nil {
		t.Errorf("Error getting payment record after locked update: %v", err)
	}
	if retrievedRecord.PaidColdRent != originalPaidColdRent {
		t.Errorf("Locked payment record values should not change")
	}

	// Test unlocking a record
	record.IsLocked = false
	err = repo.Update(record)
	if err != nil {
		t.Errorf("Error unlocking payment record: %v", err)
	}

	// Verify the record is unlocked
	retrievedRecord, err = repo.GetByID(record.ID)
	if err != nil {
		t.Errorf("Error getting unlocked payment record: %v", err)
	}
	if retrievedRecord.IsLocked {
		t.Errorf("Payment record is still locked")
	}

	// Now we should be able to update the values
	record.PaidColdRent = 900
	err = repo.Update(record)
	if err != nil {
		t.Errorf("Error updating unlocked payment record: %v", err)
	}

	// Verify the update worked
	retrievedRecord, err = repo.GetByID(record.ID)
	if err != nil {
		t.Errorf("Error getting payment record after unlock and update: %v", err)
	}
	if retrievedRecord.PaidColdRent != 900 {
		t.Errorf("Expected paid cold rent to be 900, got %f", retrievedRecord.PaidColdRent)
	}

	// Test updating with invalid data (negative payment)
	record.PaidColdRent = -100
	err = repo.Update(record)
	if err == nil {
		t.Error("Expected error for negative payment amount, got nil")
	}

	// Test updating non-existent record
	nonExistentRecord := models.PaymentRecord{
		ID:              9999,
		TenantID:        tenant.ID,
		Month:           "2023-02",
		TargetColdRent:  tenant.TargetColdRent,
		PaidColdRent:    800,
		PaidAncillary:   200,
		PaidElectricity: 100,
		ExtraPayments:   50,
		Persons:         2,
		Note:            "Non-existent record",
		IsLocked:        false,
	}

	err = repo.Update(&nonExistentRecord)
	if err == nil {
		t.Error("Expected error for updating non-existent payment record, got nil")
	}
}

func TestPaymentRecordRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDBWithPaymentRecords(t)
	defer cleanup()

	tenant := createTestTenant(t, db)
	repo := NewPaymentRecordRepository(db)

	// Create a test payment record
	record, err := models.NewPaymentRecord(
		tenant.ID,
		"2023-01",
		tenant.TargetColdRent,
		"800",
		"200",
		"100",
		"50",
		"2",
		"Payment to be deleted",
		false,
	)
	if err != nil {
		t.Fatalf("Failed to create payment record model: %v", err)
	}

	err = repo.Create(record)
	if err != nil {
		t.Fatalf("Error creating test payment record: %v", err)
	}

	// Delete the payment record
	err = repo.Delete(record.ID)
	if err != nil {
		t.Errorf("Error deleting payment record: %v", err)
	}

	// Verify the deletion
	_, err = repo.GetByID(record.ID)
	if err == nil {
		t.Error("Expected error when getting deleted payment record, got nil")
	}

	// Test deleting non-existent record
	err = repo.Delete(9999)
	if err == nil {
		t.Error("Expected error for deleting non-existent payment record, got nil")
	}
}

func TestPaymentRecordRepository_BatchCreateOrUpdateRecords(t *testing.T) {
	db, cleanup := setupTestDBWithPaymentRecords(t)
	defer cleanup()

	tenant := createTestTenant(t, db)
	repo := NewPaymentRecordRepository(db)

	// Create a batch of payment records
	var records []models.PaymentRecord
	months := []string{"2023-01", "2023-02", "2023-03"}

	for _, month := range months {
		record, err := models.NewPaymentRecord(
			tenant.ID,
			month,
			tenant.TargetColdRent,
			"800",
			"200",
			"100",
			"50",
			"2",
			"Batch payment for "+month,
			false,
		)
		if err != nil {
			t.Fatalf("Failed to create payment record model for %s: %v", month, err)
		}
		records = append(records, *record)
	}

	// Test batch creation
	err := repo.BatchCreateOrUpdateRecords(records)
	if err != nil {
		t.Errorf("Error batch creating payment records: %v", err)
	}

	// Verify the records were created
	allRecords, err := repo.GetByTenantID(tenant.ID)
	if err != nil {
		t.Errorf("Error getting all payment records: %v", err)
	}
	if len(allRecords) != len(months) {
		t.Errorf("Expected %d payment records, got %d", len(months), len(allRecords))
	}

	// Test batch update
	for i := range records {
		records[i].PaidColdRent = 850
		records[i].Note = "Updated batch payment"
	}

	err = repo.BatchCreateOrUpdateRecords(records)
	if err != nil {
		t.Errorf("Error batch updating payment records: %v", err)
	}

	// Verify the updates
	updatedRecords, err := repo.GetByTenantID(tenant.ID)
	if err != nil {
		t.Errorf("Error getting updated payment records: %v", err)
	}

	for _, record := range updatedRecords {
		if record.PaidColdRent != 850 || record.Note != "Updated batch payment" {
			t.Errorf("Payment record was not properly batch updated")
		}
	}
}

func TestPaymentRecordRepository_GeneratePaymentRecordsForTenant(t *testing.T) {
	db, cleanup := setupTestDBWithPaymentRecords(t)
	defer cleanup()

	tenant := createTestTenant(t, db)
	repo := NewPaymentRecordRepository(db)

	// Generate payment records
	err := repo.GeneratePaymentRecordsForTenant(tenant.ID)
	if err != nil {
		t.Errorf("Error generating payment records for tenant: %v", err)
	}

	// Verify records were created
	records, err := repo.GetByTenantID(tenant.ID)
	if err != nil {
		t.Errorf("Error getting generated payment records: %v", err)
	}

	// Calculate how many months from tenant's move-in date to current month
	moveInDate, _ := time.Parse("2006-01-02", "2023-01-01") // From tenant creation
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	startMonth := time.Date(moveInDate.Year(), moveInDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	expectedMonths := 0
	for m := startMonth; !m.After(currentMonth); m = m.AddDate(0, 1, 0) {
		expectedMonths++
	}

	if len(records) != expectedMonths {
		t.Errorf("Expected %d payment records, got %d", expectedMonths, len(records))
	}

	// Check if the records have the right target rent and persons
	for _, record := range records {
		if record.TargetColdRent != tenant.TargetColdRent {
			t.Errorf("Expected target cold rent %f, got %f", tenant.TargetColdRent, record.TargetColdRent)
		}
		if record.Persons != tenant.NumberOfPersons {
			t.Errorf("Expected %d persons, got %d", tenant.NumberOfPersons, record.Persons)
		}
	}

	// Test with former tenant
	formerTenant := createFormerTenant(t, db)
	err = repo.GeneratePaymentRecordsForTenant(formerTenant.ID)
	if err != nil {
		t.Errorf("Error generating payment records for former tenant: %v", err)
	}

	// Verify records for former tenant
	formerRecords, err := repo.GetByTenantID(formerTenant.ID)
	if err != nil {
		t.Errorf("Error getting generated payment records for former tenant: %v", err)
	}

	// Calculate expected months for former tenant (from move-in to move-out)
	formerMoveInDate, _ := time.Parse("2006-01-02", "2022-01-01")
	formerMoveOutDate, _ := time.Parse("2006-01-02", "2023-01-01")
	formerStartMonth := time.Date(formerMoveInDate.Year(), formerMoveInDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	formerEndMonth := time.Date(formerMoveOutDate.Year(), formerMoveOutDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	expectedFormerMonths := 0
	for m := formerStartMonth; !m.After(formerEndMonth); m = m.AddDate(0, 1, 0) {
		expectedFormerMonths++
	}

	if len(formerRecords) != expectedFormerMonths {
		t.Errorf("Expected %d payment records for former tenant, got %d", expectedFormerMonths, len(formerRecords))
	}
}

func TestGetLast12MonthsFromToday(t *testing.T) {
	months := GetLast12MonthsFromToday()

	// Verify we have 12 months
	if len(months) != 12 {
		t.Errorf("Expected 12 months, got %d", len(months))
	}

	// Verify format and order
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Last month in the list should be current month
	lastMonth := months[len(months)-1]
	expectedLastMonth := currentMonth.Format("2006-01")
	if lastMonth != expectedLastMonth {
		t.Errorf("Expected last month to be %s, got %s", expectedLastMonth, lastMonth)
	}

	// First month in the list should be 11 months ago
	firstMonth := months[0]
	expectedFirstMonth := currentMonth.AddDate(0, -11, 0).Format("2006-01")
	if firstMonth != expectedFirstMonth {
		t.Errorf("Expected first month to be %s, got %s", expectedFirstMonth, firstMonth)
	}

	// Verify the months are in ascending order
	for i := 1; i < len(months); i++ {
		if months[i] <= months[i-1] {
			t.Errorf("Months are not in ascending order: %s after %s", months[i], months[i-1])
		}
	}
}

func TestGetMonthsRange(t *testing.T) {
	// Test valid range
	months, err := GetMonthsRange("2023-01", "2023-06")
	if err != nil {
		t.Errorf("Error getting months range: %v", err)
	}
	if len(months) != 6 {
		t.Errorf("Expected 6 months, got %d", len(months))
	}
	if months[0] != "2023-01" || months[5] != "2023-06" {
		t.Errorf("Incorrect months range: %v", months)
	}

	// Test range spanning multiple years
	months, err = GetMonthsRange("2022-11", "2023-02")
	if err != nil {
		t.Errorf("Error getting months range spanning years: %v", err)
	}
	if len(months) != 4 {
		t.Errorf("Expected 4 months, got %d", len(months))
	}
	if months[0] != "2022-11" || months[3] != "2023-02" {
		t.Errorf("Incorrect months range spanning years: %v", months)
	}

	// Test invalid start month
	_, err = GetMonthsRange("2023-13", "2023-06")
	if err == nil {
		t.Errorf("Expected error for invalid start month, got nil")
	}

	// Test invalid end month
	_, err = GetMonthsRange("2023-01", "2023-13")
	if err == nil {
		t.Errorf("Expected error for invalid end month, got nil")
	}

	// Test end before start
	_, err = GetMonthsRange("2023-06", "2023-01")
	if err == nil {
		t.Errorf("Expected error for end month before start month, got nil")
	}
}
