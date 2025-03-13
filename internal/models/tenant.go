package models

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"property-management/internal/utils"
)

// Tenant represents a tenant living in an apartment
type Tenant struct {
	ID                       int64      `json:"id"`
	FirstName                string     `json:"firstName"`
	LastName                 string     `json:"lastName"`
	MoveInDate               time.Time  `json:"moveInDate"`
	MoveOutDate              *time.Time `json:"moveOutDate,omitempty"`
	Deposit                  float64    `json:"deposit"`
	Email                    string     `json:"email,omitempty"`
	NumberOfPersons          int        `json:"numberOfPersons"`
	TargetColdRent           float64    `json:"targetColdRent"`
	TargetAncillaryPayment   float64    `json:"targetAncillaryPayment"`
	TargetElectricityPayment float64    `json:"targetElectricityPayment"`
	Greeting                 string     `json:"greeting"`
	HouseID                  int64      `json:"houseId"`
	House                    *House     `json:"house,omitempty"`
	ApartmentID              int64      `json:"apartmentId"`
	Apartment                *Apartment `json:"apartment,omitempty"`
	CreatedAt                time.Time  `json:"createdAt"`
	UpdatedAt                time.Time  `json:"updatedAt"`
}

// NewTenant creates a new tenant with the given details
func NewTenant(
	firstName string,
	lastName string,
	moveInDateStr string,
	moveOutDateStr string,
	depositStr string,
	email string,
	numberOfPersonsStr string,
	targetColdRentStr string,
	targetAncillaryPaymentStr string,
	targetElectricityPaymentStr string,
	greeting string,
	houseID int64,
	apartmentID int64,
) (*Tenant, error) {
	// Parse move-in date
	moveInDate, err := parseDate(moveInDateStr)
	if err != nil {
		return nil, errors.New("invalid move-in date format: must be YYYY-MM-DD")
	}

	// Parse optional move-out date
	var moveOutDate *time.Time
	if moveOutDateStr != "" {
		parsed, err := parseDate(moveOutDateStr)
		if err != nil {
			return nil, errors.New("invalid move-out date format: must be YYYY-MM-DD")
		}
		moveOutDate = &parsed

		// Check if move-out date is after move-in date
		if !moveOutDate.After(moveInDate) {
			return nil, errors.New("move-out date must be after move-in date")
		}
	}

	// Parse deposit
	deposit, err := utils.ParseAmount(depositStr)
	if err != nil {
		return nil, errors.New("invalid deposit amount")
	}

	// Parse number of persons
	numberOfPersons, err := strconv.Atoi(strings.TrimSpace(numberOfPersonsStr))
	if err != nil {
		return nil, errors.New("invalid number of persons")
	}
	if numberOfPersons <= 0 {
		return nil, errors.New("number of persons must be greater than 0")
	}

	// Parse target cold rent
	targetColdRent, err := utils.ParseAmount(targetColdRentStr)
	if err != nil {
		return nil, errors.New("invalid target cold rent amount")
	}
	if targetColdRent <= 0 {
		return nil, errors.New("target cold rent must be greater than 0")
	}

	// Parse target ancillary payment
	targetAncillaryPayment, err := utils.ParseAmount(targetAncillaryPaymentStr)
	if err != nil {
		return nil, errors.New("invalid target ancillary payment amount")
	}
	if targetAncillaryPayment <= 0 {
		return nil, errors.New("target ancillary payment must be greater than 0")
	}

	// Parse target electricity payment
	targetElectricityPayment, err := utils.ParseAmount(targetElectricityPaymentStr)
	if err != nil {
		return nil, errors.New("invalid target electricity payment amount")
	}
	if targetElectricityPayment <= 0 {
		return nil, errors.New("target electricity payment must be greater than 0")
	}

	// Check HouseID and ApartmentID
	if houseID <= 0 {
		return nil, errors.New("invalid house ID")
	}

	if apartmentID <= 0 {
		return nil, errors.New("invalid apartment ID")
	}

	now := time.Now()
	return &Tenant{
		FirstName:                firstName,
		LastName:                 lastName,
		MoveInDate:               moveInDate,
		MoveOutDate:              moveOutDate,
		Deposit:                  deposit,
		Email:                    email,
		NumberOfPersons:          numberOfPersons,
		TargetColdRent:           targetColdRent,
		TargetAncillaryPayment:   targetAncillaryPayment,
		TargetElectricityPayment: targetElectricityPayment,
		Greeting:                 greeting,
		HouseID:                  houseID,
		ApartmentID:              apartmentID,
		CreatedAt:                now,
		UpdatedAt:                now,
	}, nil
}

// Validate ensures all tenant data is valid
func (t *Tenant) Validate() error {
	// First and Last name validation
	if strings.TrimSpace(t.FirstName) == "" {
		return errors.New("first name cannot be empty")
	}

	if strings.TrimSpace(t.LastName) == "" {
		return errors.New("last name cannot be empty")
	}

	// MoveInDate validation is implicit (time.Time can't be "empty")

	// If MoveOutDate is specified, it should be after MoveInDate
	if t.MoveOutDate != nil && !t.MoveOutDate.After(t.MoveInDate) {
		return errors.New("move-out date must be after move-in date")
	}

	// NumberOfPersons validation
	if t.NumberOfPersons <= 0 {
		return errors.New("number of persons must be greater than 0")
	}

	// Payment amounts validation
	if t.TargetColdRent <= 0 {
		return errors.New("target cold rent must be greater than 0")
	}

	if t.TargetAncillaryPayment <= 0 {
		return errors.New("target ancillary payment must be greater than 0")
	}

	if t.TargetElectricityPayment <= 0 {
		return errors.New("target electricity payment must be greater than 0")
	}

	// Deposit amount allows zero (no deposit)
	if t.Deposit < 0 {
		return errors.New("deposit cannot be negative")
	}

	// House and Apartment IDs validation
	if t.HouseID <= 0 {
		return errors.New("invalid house ID")
	}

	if t.ApartmentID <= 0 {
		return errors.New("invalid apartment ID")
	}

	return nil
}

// Helper functions
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(dateStr))
}
