package models

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"property-management/internal/utils"
)

// PaymentRecord represents a monthly payment record for a tenant
type PaymentRecord struct {
	ID              int64     `json:"id"`
	TenantID        int64     `json:"tenantId"`
	Tenant          *Tenant   `json:"tenant,omitempty"`
	Month           string    `json:"month"` // Format: YYYY-MM
	TargetColdRent  float64   `json:"targetColdRent"`
	PaidColdRent    float64   `json:"paidColdRent"`
	PaidAncillary   float64   `json:"paidAncillary"`
	PaidElectricity float64   `json:"paidElectricity"`
	ExtraPayments   float64   `json:"extraPayments"`
	Persons         int       `json:"persons"`
	Note            string    `json:"note"`
	IsLocked        bool      `json:"isLocked"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// NewPaymentRecord creates a new payment record
func NewPaymentRecord(
	tenantID int64,
	month string,
	targetColdRent float64,
	paidColdRentStr string,
	paidAncillaryStr string,
	paidElectricityStr string,
	extraPaymentsStr string,
	personsStr string,
	note string,
	isLocked bool,
) (*PaymentRecord, error) {
	// Validate month format (YYYY-MM)
	if _, err := time.Parse("2006-01", month); err != nil {
		return nil, errors.New("invalid month format: must be YYYY-MM")
	}

	// Parse payment values
	paidColdRent, err := utils.ParseAmount(paidColdRentStr)
	if err != nil {
		return nil, errors.New("invalid paid cold rent amount")
	}

	paidAncillary, err := utils.ParseAmount(paidAncillaryStr)
	if err != nil {
		return nil, errors.New("invalid paid ancillary amount")
	}

	paidElectricity, err := utils.ParseAmount(paidElectricityStr)
	if err != nil {
		return nil, errors.New("invalid paid electricity amount")
	}

	extraPayments, err := utils.ParseAmount(extraPaymentsStr)
	if err != nil {
		return nil, errors.New("invalid extra payments amount")
	}

	// Parse persons
	persons, err := strconv.Atoi(strings.TrimSpace(personsStr))
	if err != nil {
		return nil, errors.New("invalid number of persons")
	}
	if persons <= 0 {
		return nil, errors.New("number of persons must be greater than 0")
	}

	now := time.Now()
	return &PaymentRecord{
		TenantID:        tenantID,
		Month:           month,
		TargetColdRent:  targetColdRent,
		PaidColdRent:    paidColdRent,
		PaidAncillary:   paidAncillary,
		PaidElectricity: paidElectricity,
		ExtraPayments:   extraPayments,
		Persons:         persons,
		Note:            note,
		IsLocked:        isLocked,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Validate ensures all payment record data is valid
func (p *PaymentRecord) Validate() error {
	// Validate TenantID
	if p.TenantID <= 0 {
		return errors.New("invalid tenant ID")
	}

	// Validate month format (YYYY-MM)
	if _, err := time.Parse("2006-01", p.Month); err != nil {
		return errors.New("invalid month format: must be YYYY-MM")
	}

	// Validate persons
	if p.Persons <= 0 {
		return errors.New("number of persons must be greater than 0")
	}

	// All amounts should be non-negative
	if p.PaidColdRent < 0 || p.PaidAncillary < 0 || p.PaidElectricity < 0 || p.ExtraPayments < 0 {
		return errors.New("payment amounts cannot be negative")
	}

	return nil
}

// FormatMonth formats the month string to a human-readable format (MM.YYYY)
func (p *PaymentRecord) FormatMonth() string {
	t, err := time.Parse("2006-01", p.Month)
	if err != nil {
		return p.Month // Return original if parsing fails
	}
	return t.Format("01.2006") // MM.YYYY format
}

// DisplayMonth formats the month string to a human-readable format (e.g., "September 2024")
func (p *PaymentRecord) DisplayMonth() string {
	t, err := time.Parse("2006-01", p.Month)
	if err != nil {
		return p.Month // Return original if parsing fails
	}
	return t.Format("January 2006")
}
