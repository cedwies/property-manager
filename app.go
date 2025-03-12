package main

import (
	"context"
	"database/sql"

	"property-management/internal/db"
	"property-management/internal/models"
	"property-management/internal/repository"
)

// App struct represents the application
type App struct {
	ctx             context.Context
	db              *sql.DB
	houseRepository *repository.HouseRepository
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
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	db.Close()
}

// GetAppInfo returns basic information about the application
func (a *App) GetAppInfo() map[string]string {
	return map[string]string{
		"name":    "Property Management System",
		"version": "0.1.0",
		"status":  "Initial Setup - Houses Management Implemented",
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
