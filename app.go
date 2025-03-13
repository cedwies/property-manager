package main

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"property-management/internal/db"
	"property-management/internal/models"
	"property-management/internal/repository"
	"property-management/internal/utils"
)

// App struct represents the application
type App struct {
	ctx                     context.Context
	db                      *sql.DB
	houseRepository         *repository.HouseRepository
	apartmentRepository     *repository.ApartmentRepository
	tenantRepository        *repository.TenantRepository
	paymentRecordRepository *repository.PaymentRecordRepository
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.db = db.GetDB()
	a.houseRepository = repository.NewHouseRepository(a.db)
	a.apartmentRepository = repository.NewApartmentRepository(a.db)
	a.tenantRepository = repository.NewTenantRepository(a.db)
	a.paymentRecordRepository = repository.NewPaymentRecordRepository(a.db)

	// Ensure payment records exist for active tenants
	err := a.paymentRecordRepository.EnsurePaymentRecordsForActiveTenants()
	if err != nil {
		// Log error but don't crash the application
		println("Error ensuring payment records for active tenants:", err.Error())
	}
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	db.Close()
}

// GetAppInfo returns basic information about the application
func (a *App) GetAppInfo() map[string]string {
	return map[string]string{
		"name":    "Property Management System",
		"version": "0.4.0",
		"status":  "Houses, Apartments, Tenants, and Payments Management Implemented",
	}
}

// CreateHouse adds a new house to the database
func (a *App) CreateHouse(name, street, number, country, zipCode, city string) (*models.House, error) {
	house := models.NewHouse(name, street, number, country, zipCode, city)
	err := a.houseRepository.Create(house)
	if err != nil {
		return nil, err
	}
	return house, nil
}

// GetAllHouses returns all houses from the database
func (a *App) GetAllHouses() ([]models.House, error) {
	return a.houseRepository.GetAll()
}

// GetHouseByID returns a house with the specified ID
func (a *App) GetHouseByID(id int64) (*models.House, error) {
	return a.houseRepository.GetByID(id)
}

// UpdateHouse modifies an existing house in the database
func (a *App) UpdateHouse(id int64, name, street, number, country, zipCode, city string) (*models.House, error) {
	house := &models.House{
		ID:      id,
		Name:    name,
		Street:  street,
		Number:  number,
		Country: country,
		ZipCode: zipCode,
		City:    city,
	}

	err := a.houseRepository.Update(house)
	if err != nil {
		return nil, err
	}

	return house, nil
}

// DeleteHouse removes a house from the database
func (a *App) DeleteHouse(id int64) error {
	return a.houseRepository.Delete(id)
}

// CreateApartment adds a new apartment to the database
func (a *App) CreateApartment(name string, houseID int64, size string) (*models.Apartment, error) {
	apartment, err := models.NewApartment(name, houseID, size)
	if err != nil {
		return nil, err
	}

	err = a.apartmentRepository.Create(apartment)
	if err != nil {
		return nil, err
	}

	// Fetch the complete apartment with house information
	return a.apartmentRepository.GetByID(apartment.ID)
}

// GetAllApartments returns all apartments from the database
func (a *App) GetAllApartments() ([]models.Apartment, error) {
	return a.apartmentRepository.GetAll()
}

// GetApartmentsByHouseID returns all apartments for a specific house
func (a *App) GetApartmentsByHouseID(houseID int64) ([]models.Apartment, error) {
	return a.apartmentRepository.GetByHouseID(houseID)
}

// GetApartmentByID returns an apartment with the specified ID
func (a *App) GetApartmentByID(id int64) (*models.Apartment, error) {
	return a.apartmentRepository.GetByID(id)
}

// UpdateApartment modifies an existing apartment in the database
func (a *App) UpdateApartment(id int64, name string, houseID int64, size string) (*models.Apartment, error) {
	// Parse size, supporting both dot and comma as decimal separators
	size = strings.TrimSpace(size)
	// Replace comma with dot for parsing
	size = strings.Replace(size, ",", ".", -1)

	sizeFloat, err := strconv.ParseFloat(size, 64)
	if err != nil {
		return nil, errors.New("invalid size format")
	}

	if sizeFloat <= 0 {
		return nil, errors.New("size must be greater than 0")
	}

	apartment := &models.Apartment{
		ID:      id,
		Name:    name,
		HouseID: houseID,
		Size:    sizeFloat,
	}

	err = a.apartmentRepository.Update(apartment)
	if err != nil {
		return nil, err
	}

	// Fetch the updated apartment with house information
	return a.apartmentRepository.GetByID(id)
}

// DeleteApartment removes an apartment from the database
func (a *App) DeleteApartment(id int64) error {
	return a.apartmentRepository.Delete(id)
}

// CreateTenant adds a new tenant to the database
func (a *App) CreateTenant(firstName, lastName, moveInDate, moveOutDate, deposit,
	email, numberOfPersons, targetColdRent, targetAncillaryPayment,
	targetElectricityPayment, greeting string, houseID, apartmentID int64) (*models.Tenant, error) {

	tenant, err := models.NewTenant(
		firstName,
		lastName,
		moveInDate,
		moveOutDate,
		deposit,
		email,
		numberOfPersons,
		targetColdRent,
		targetAncillaryPayment,
		targetElectricityPayment,
		greeting,
		houseID,
		apartmentID,
	)
	if err != nil {
		return nil, err
	}

	err = a.tenantRepository.Create(tenant)
	if err != nil {
		return nil, err
	}

	// Generate payment records for the new tenant
	err = a.paymentRecordRepository.GeneratePaymentRecordsForTenant(tenant.ID)
	if err != nil {
		// Log error but don't fail the tenant creation
		println("Error generating payment records for new tenant:", err.Error())
	}

	// Fetch the complete tenant with house and apartment information
	return a.tenantRepository.GetByID(tenant.ID)
}

// GetAllTenants returns all tenants from the database
func (a *App) GetAllTenants() ([]models.Tenant, error) {
	return a.tenantRepository.GetAll()
}

// GetTenantsByHouseID returns all tenants for a specific house
func (a *App) GetTenantsByHouseID(houseID int64) ([]models.Tenant, error) {
	return a.tenantRepository.GetByHouseID(houseID)
}

// GetCurrentTenantsByHouseID returns all current tenants for a specific house
func (a *App) GetCurrentTenantsByHouseID(houseID int64) ([]models.Tenant, error) {
	allTenants, err := a.tenantRepository.GetByHouseID(houseID)
	if err != nil {
		return nil, err
	}

	currentTenants := []models.Tenant{}
	now := time.Now()

	for _, tenant := range allTenants {
		// If tenant has no move-out date or move-out date is in the future
		if tenant.MoveOutDate == nil || tenant.MoveOutDate.After(now) {
			currentTenants = append(currentTenants, tenant)
		}
	}

	return currentTenants, nil
}

// GetTenantsByApartmentID returns all tenants for a specific apartment
func (a *App) GetTenantsByApartmentID(apartmentID int64) ([]models.Tenant, error) {
	return a.tenantRepository.GetByApartmentID(apartmentID)
}

// GetTenantByID returns a tenant with the specified ID
func (a *App) GetTenantByID(id int64) (*models.Tenant, error) {
	return a.tenantRepository.GetByID(id)
}

// UpdateTenant modifies an existing tenant in the database
func (a *App) UpdateTenant(
	id int64,
	firstName,
	lastName,
	moveInDate,
	moveOutDate,
	deposit,
	email,
	numberOfPersons,
	targetColdRent,
	targetAncillaryPayment,
	targetElectricityPayment,
	greeting string,
	houseID,
	apartmentID int64) (*models.Tenant, error) {

	// First fetch the existing tenant to preserve any fields we don't want to modify
	existingTenant, err := a.tenantRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Create a new tenant object with the updated information
	updatedTenant, err := models.NewTenant(
		firstName,
		lastName,
		moveInDate,
		moveOutDate,
		deposit,
		email,
		numberOfPersons,
		targetColdRent,
		targetAncillaryPayment,
		targetElectricityPayment,
		greeting,
		houseID,
		apartmentID,
	)
	if err != nil {
		return nil, err
	}

	// Set the ID and timestamps
	updatedTenant.ID = id
	updatedTenant.CreatedAt = existingTenant.CreatedAt

	// Update the tenant in the database
	err = a.tenantRepository.Update(updatedTenant)
	if err != nil {
		return nil, err
	}

	// If target cold rent changed, we may need to update payment records
	if existingTenant.TargetColdRent != updatedTenant.TargetColdRent {
		// This is a simplified approach - in a real system, you might want to
		// update only future payment records or ask the user which records to update
		err = a.updateFuturePaymentRecordTargetRents(id, updatedTenant.TargetColdRent)
		if err != nil {
			// Log error but don't fail the tenant update
			println("Error updating payment records target rents:", err.Error())
		}
	}

	// If the number of persons changed, we may need to update payment records
	if existingTenant.NumberOfPersons != updatedTenant.NumberOfPersons {
		err = a.updateFuturePaymentRecordPersons(id, updatedTenant.NumberOfPersons)
		if err != nil {
			// Log error but don't fail the tenant update
			println("Error updating payment records persons:", err.Error())
		}
	}

	// Fetch the complete updated tenant with house and apartment information
	return a.tenantRepository.GetByID(id)
}

// DeleteTenant removes a tenant from the database
func (a *App) DeleteTenant(id int64) error {
	return a.tenantRepository.Delete(id)
}

// Payment Record Methods

// CreatePaymentRecord adds a new payment record to the database
func (a *App) CreatePaymentRecord(
	tenantID int64,
	month string,
	targetColdRent float64,
	paidColdRent string,
	paidAncillary string,
	paidElectricity string,
	extraPayments string,
	persons string,
	note string,
	isLocked bool) (*models.PaymentRecord, error) {

	record, err := models.NewPaymentRecord(
		tenantID,
		month,
		targetColdRent,
		paidColdRent,
		paidAncillary,
		paidElectricity,
		extraPayments,
		persons,
		note,
		isLocked,
	)
	if err != nil {
		return nil, err
	}

	err = a.paymentRecordRepository.Create(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

// GetPaymentRecordByID returns a payment record with the specified ID
func (a *App) GetPaymentRecordByID(id int64) (*models.PaymentRecord, error) {
	return a.paymentRecordRepository.GetByID(id)
}

// GetPaymentRecordsByTenantID returns all payment records for a specific tenant
func (a *App) GetPaymentRecordsByTenantID(tenantID int64) ([]models.PaymentRecord, error) {
	return a.paymentRecordRepository.GetByTenantID(tenantID)
}

// GetRecentPaymentRecordsByTenantID returns recent payment records (last 12 months) for a specific tenant
func (a *App) GetRecentPaymentRecordsByTenantID(tenantID int64) ([]models.PaymentRecord, error) {
	return a.paymentRecordRepository.GetRecentByTenantID(tenantID, 12)
}

// GetPaymentRecordByTenantAndMonth returns a payment record for a specific tenant and month
func (a *App) GetPaymentRecordByTenantAndMonth(tenantID int64, month string) (*models.PaymentRecord, error) {
	return a.paymentRecordRepository.GetByTenantAndMonth(tenantID, month)
}

// UpdatePaymentRecord modifies an existing payment record in the database
func (a *App) UpdatePaymentRecord(
	id int64,
	paidColdRent string,
	paidAncillary string,
	paidElectricity string,
	extraPayments string,
	persons string,
	note string,
	isLocked bool) (*models.PaymentRecord, error) {

	// First fetch the existing record
	existingRecord, err := a.paymentRecordRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Parse payment values
	paidColdRentVal, err := utils.ParseAmount(paidColdRent)
	if err != nil {
		return nil, errors.New("invalid paid cold rent amount")
	}

	paidAncillaryVal, err := utils.ParseAmount(paidAncillary)
	if err != nil {
		return nil, errors.New("invalid paid ancillary amount")
	}

	paidElectricityVal, err := utils.ParseAmount(paidElectricity)
	if err != nil {
		return nil, errors.New("invalid paid electricity amount")
	}

	extraPaymentsVal, err := utils.ParseAmount(extraPayments)
	if err != nil {
		return nil, errors.New("invalid extra payments amount")
	}

	// Parse persons
	personsVal, err := strconv.Atoi(strings.TrimSpace(persons))
	if err != nil {
		return nil, errors.New("invalid number of persons")
	}
	if personsVal <= 0 {
		return nil, errors.New("number of persons must be greater than 0")
	}

	// Update the record
	existingRecord.PaidColdRent = paidColdRentVal
	existingRecord.PaidAncillary = paidAncillaryVal
	existingRecord.PaidElectricity = paidElectricityVal
	existingRecord.ExtraPayments = extraPaymentsVal
	existingRecord.Persons = personsVal
	existingRecord.Note = note
	existingRecord.IsLocked = isLocked

	err = a.paymentRecordRepository.Update(existingRecord)
	if err != nil {
		return nil, err
	}

	return existingRecord, nil
}

// DeletePaymentRecord removes a payment record from the database
func (a *App) DeletePaymentRecord(id int64) error {
	return a.paymentRecordRepository.Delete(id)
}

// GetPaymentRecordsForHouse returns all payment records for current tenants in a house
func (a *App) GetPaymentRecordsForHouse(houseID int64) (map[int64][]models.PaymentRecord, error) {
	return a.paymentRecordRepository.GetCurrentTenantsPaymentsByHouseID(houseID)
}

// BatchSavePaymentRecords updates or creates multiple payment records at once
func (a *App) BatchSavePaymentRecords(records []models.PaymentRecord) error {
	return a.paymentRecordRepository.BatchCreateOrUpdateRecords(records)
}

// GetLastTwelveMonths returns the last 12 months for payment record display
func (a *App) GetLastTwelveMonths() []string {
	return repository.GetLast12MonthsFromToday()
}

// UpdateTenantPaymentRecordsNote updates the note for a payment record
func (a *App) UpdatePaymentRecordNote(recordID int64, note string) error {
	record, err := a.paymentRecordRepository.GetByID(recordID)
	if err != nil {
		return err
	}

	record.Note = note
	return a.paymentRecordRepository.Update(record)
}

// TogglePaymentRecordLock toggles the lock state of a payment record
func (a *App) TogglePaymentRecordLock(recordID int64) error {
	record, err := a.paymentRecordRepository.GetByID(recordID)
	if err != nil {
		return err
	}

	record.IsLocked = !record.IsLocked
	return a.paymentRecordRepository.Update(record)
}

// updateFuturePaymentRecordTargetRents updates the target rent for future payment records
func (a *App) updateFuturePaymentRecordTargetRents(tenantID int64, newTargetRent float64) error {
	// Get all payment records for this tenant
	records, err := a.paymentRecordRepository.GetByTenantID(tenantID)
	if err != nil {
		return err
	}

	// Get current month
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	currentMonthStr := currentMonth.Format("2006-01")

	// Update records for current month and future months
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}

	for _, record := range records {
		// Skip records for past months
		if record.Month < currentMonthStr {
			continue
		}

		// Skip locked records
		if record.IsLocked {
			continue
		}

		// Update target rent
		query := "UPDATE payment_records SET target_cold_rent = ? WHERE id = ?"
		_, err := tx.Exec(query, newTargetRent, record.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// updateFuturePaymentRecordPersons updates the number of persons for future payment records
func (a *App) updateFuturePaymentRecordPersons(tenantID int64, newPersons int) error {
	// Get all payment records for this tenant
	records, err := a.paymentRecordRepository.GetByTenantID(tenantID)
	if err != nil {
		return err
	}

	// Get current month
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	currentMonthStr := currentMonth.Format("2006-01")

	// Update records for current month and future months
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}

	for _, record := range records {
		// Skip records for past months
		if record.Month < currentMonthStr {
			continue
		}

		// Skip locked records
		if record.IsLocked {
			continue
		}

		// Update persons
		query := "UPDATE payment_records SET persons = ? WHERE id = ?"
		_, err := tx.Exec(query, newPersons, record.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
