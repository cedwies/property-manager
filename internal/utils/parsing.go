package utils

import (
	"errors"
	"strconv"
	"strings"
)

// ParseAmount parses a string amount, supporting both dot and comma as decimal separators
// Returns the parsed float64 value or an error if the amount is invalid or negative
func ParseAmount(amountStr string) (float64, error) {
	if amountStr == "" {
		return 0, nil // Allow empty string for optional amounts
	}

	// Support both dot and comma as decimal separators
	amountStr = strings.TrimSpace(amountStr)
	// Replace comma with dot for parsing
	amountStr = strings.Replace(amountStr, ",", ".", -1)

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return 0, err
	}

	if amount < 0 {
		return 0, errors.New("amount cannot be negative")
	}

	return amount, nil
}
