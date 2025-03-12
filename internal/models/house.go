package models

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// House represents a property in the system
type House struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Street    string    `json:"street"`
	Number    string    `json:"number"`
	Country   string    `json:"country"`
	ZipCode   string    `json:"zipCode"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Validate ensures all house data is valid
func (h *House) Validate() error {
	// Name validation
	if strings.TrimSpace(h.Name) == "" {
		return errors.New("house name cannot be empty")
	}

	// Street validation
	if strings.TrimSpace(h.Street) == "" {
		return errors.New("street cannot be empty")
	}

	// Number validation
	if strings.TrimSpace(h.Number) == "" {
		return errors.New("house number cannot be empty")
	}

	// Country validation
	if strings.TrimSpace(h.Country) == "" {
		return errors.New("country cannot be empty")
	}

	// ZipCode validation
	if strings.TrimSpace(h.ZipCode) == "" {
		return errors.New("zip code cannot be empty")
	}

	// City validation - should not be just a number
	city := strings.TrimSpace(h.City)
	if city == "" {
		return errors.New("city cannot be empty")
	}

	// Check if city is only numeric
	re := regexp.MustCompile(`^\d+$`)
	if re.MatchString(city) {
		return errors.New("city cannot be just a number")
	}

	return nil
}

// NewHouse creates a new house with the given details
func NewHouse(name, street, number, country, zipCode, city string) *House {
	now := time.Now()
	return &House{
		Name:      name,
		Street:    street,
		Number:    number,
		Country:   country,
		ZipCode:   zipCode,
		City:      city,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
