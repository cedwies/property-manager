package models

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Apartment represents an apartment within a house
type Apartment struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	HouseID   int64     `json:"houseId"`
	House     *House    `json:"house,omitempty"` // The house this apartment belongs to
	Size      float64   `json:"size"`            // Size in square meters
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewApartment creates a new apartment with the given details
func NewApartment(name string, houseID int64, sizeStr string) (*Apartment, error) {
	// Parse size, supporting both dot and comma as decimal separators
	sizeStr = strings.TrimSpace(sizeStr)
	// Replace comma with dot for parsing
	sizeStr = strings.Replace(sizeStr, ",", ".", -1)

	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return nil, errors.New("invalid size format")
	}

	if size <= 0 {
		return nil, errors.New("size must be greater than 0")
	}

	now := time.Now()
	return &Apartment{
		Name:      name,
		HouseID:   houseID,
		Size:      size,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Validate ensures all apartment data is valid
func (a *Apartment) Validate() error {
	// Name validation
	if strings.TrimSpace(a.Name) == "" {
		return errors.New("apartment name cannot be empty")
	}

	// House ID validation
	if a.HouseID <= 0 {
		return errors.New("invalid house ID")
	}

	// Size validation
	if a.Size <= 0 {
		return errors.New("size must be greater than 0")
	}

	return nil
}
