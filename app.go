package main

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"property-management/internal/db"
	"property-management/internal/models"
	"property-management/internal/repository"
)

// App struct represents the application
type App struct {
	ctx                 context.Context
	db                  *sql.DB
	houseRepository     *repository.HouseRepository
	apartmentRepository *repository.ApartmentRepository
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
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	db.Close()
}

// GetAppInfo returns basic information about the application
func (a *App) GetAppInfo() map[string]string {
	return map[string]string{
		"name":    "Property Management System",
		"version": "0.2.0",
		"status":  "Houses and Apartments Management Implemented",
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
