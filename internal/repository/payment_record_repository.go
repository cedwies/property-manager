package repository

import (
	"database/sql"
	"errors"
	"property-management/internal/models"
	"strings"
	"time"
)

// PaymentRecordRepository handles all database interactions for payment records
type PaymentRecordRepository struct {
	db *sql.DB
}

// NewPaymentRecordRepository creates a new payment record repository
func NewPaymentRecordRepository(db *sql.DB) *PaymentRecordRepository {
	return &PaymentRecordRepository{db: db}
}

// InitSchema initializes the payment records table schema
func (r *PaymentRecordRepository) InitSchema() error {
	// Create payment_records table
	query := `
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

	_, err := r.db.Exec(query)
	return err
}

// Create adds a new payment record to the database
func (r *PaymentRecordRepository) Create(record *models.PaymentRecord) error {
	// Validate payment record data
	if err := record.Validate(); err != nil {
		return err
	}

	// Check if tenant exists
	var tenantExists int
	err := r.db.QueryRow("SELECT COUNT(*) FROM tenants WHERE id = ?", record.TenantID).Scan(&tenantExists)
	if err != nil {
		return err
	}
	if tenantExists == 0 {
		return errors.New("tenant not found")
	}

	// Check if record for this tenant and month already exists
	var count int
	err = r.db.QueryRow("SELECT COUNT(*) FROM payment_records WHERE tenant_id = ? AND month = ?",
		record.TenantID, record.Month).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("payment record already exists for this tenant and month")
	}

	// Prepare the SQL statement
	query := `
		INSERT INTO payment_records (
			tenant_id, month, target_cold_rent, paid_cold_rent, paid_ancillary,
			paid_electricity, extra_payments, persons, note, is_locked,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Execute the query
	now := time.Now()
	result, err := r.db.Exec(
		query,
		record.TenantID,
		record.Month,
		record.TargetColdRent,
		record.PaidColdRent,
		record.PaidAncillary,
		record.PaidElectricity,
		record.ExtraPayments,
		record.Persons,
		record.Note,
		record.IsLocked,
		now,
		now,
	)
	if err != nil {
		return err
	}

	// Get the inserted ID and update the record object
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	record.ID = id
	record.CreatedAt = now
	record.UpdatedAt = now

	return nil
}

// GetByID returns a payment record with the specified ID
func (r *PaymentRecordRepository) GetByID(id int64) (*models.PaymentRecord, error) {
	// Prepare the SQL statement
	query := `
		SELECT pr.id, pr.tenant_id, pr.month, pr.target_cold_rent, pr.paid_cold_rent,
			pr.paid_ancillary, pr.paid_electricity, pr.extra_payments, pr.persons,
			pr.note, pr.is_locked, pr.created_at, pr.updated_at
		FROM payment_records pr
		WHERE pr.id = ?
	`

	// Execute the query
	var record models.PaymentRecord
	var createdAt, updatedAt string

	err := r.db.QueryRow(query, id).Scan(
		&record.ID,
		&record.TenantID,
		&record.Month,
		&record.TargetColdRent,
		&record.PaidColdRent,
		&record.PaidAncillary,
		&record.PaidElectricity,
		&record.ExtraPayments,
		&record.Persons,
		&record.Note,
		&record.IsLocked,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment record not found")
		}
		return nil, err
	}

	// Parse timestamps
	record.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	record.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &record, nil
}

// GetByTenantAndMonth returns a payment record for a specific tenant and month
func (r *PaymentRecordRepository) GetByTenantAndMonth(tenantID int64, month string) (*models.PaymentRecord, error) {
	// Prepare the SQL statement
	query := `
		SELECT pr.id, pr.tenant_id, pr.month, pr.target_cold_rent, pr.paid_cold_rent,
			pr.paid_ancillary, pr.paid_electricity, pr.extra_payments, pr.persons,
			pr.note, pr.is_locked, pr.created_at, pr.updated_at
		FROM payment_records pr
		WHERE pr.tenant_id = ? AND pr.month = ?
	`

	// Execute the query
	var record models.PaymentRecord
	var createdAt, updatedAt string

	err := r.db.QueryRow(query, tenantID, month).Scan(
		&record.ID,
		&record.TenantID,
		&record.Month,
		&record.TargetColdRent,
		&record.PaidColdRent,
		&record.PaidAncillary,
		&record.PaidElectricity,
		&record.ExtraPayments,
		&record.Persons,
		&record.Note,
		&record.IsLocked,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment record not found")
		}
		return nil, err
	}

	// Parse timestamps
	record.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	record.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &record, nil
}

// GetByTenantID returns all payment records for a specific tenant
func (r *PaymentRecordRepository) GetByTenantID(tenantID int64) ([]models.PaymentRecord, error) {
	// Prepare the SQL statement
	query := `
		SELECT pr.id, pr.tenant_id, pr.month, pr.target_cold_rent, pr.paid_cold_rent,
			pr.paid_ancillary, pr.paid_electricity, pr.extra_payments, pr.persons,
			pr.note, pr.is_locked, pr.created_at, pr.updated_at
		FROM payment_records pr
		WHERE pr.tenant_id = ?
		ORDER BY pr.month DESC
	`

	// Execute the query
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	var records []models.PaymentRecord
	for rows.Next() {
		var record models.PaymentRecord
		var createdAt, updatedAt string

		err := rows.Scan(
			&record.ID,
			&record.TenantID,
			&record.Month,
			&record.TargetColdRent,
			&record.PaidColdRent,
			&record.PaidAncillary,
			&record.PaidElectricity,
			&record.ExtraPayments,
			&record.Persons,
			&record.Note,
			&record.IsLocked,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamps
		record.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		record.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// GetRecentByTenantID returns the most recent payment records for a specific tenant
// limited by the provided count parameter
func (r *PaymentRecordRepository) GetRecentByTenantID(tenantID int64, count int) ([]models.PaymentRecord, error) {
	// Prepare the SQL statement
	query := `
		SELECT pr.id, pr.tenant_id, pr.month, pr.target_cold_rent, pr.paid_cold_rent,
			pr.paid_ancillary, pr.paid_electricity, pr.extra_payments, pr.persons,
			pr.note, pr.is_locked, pr.created_at, pr.updated_at
		FROM payment_records pr
		WHERE pr.tenant_id = ?
		ORDER BY pr.month DESC
		LIMIT ?
	`

	// Execute the query
	rows, err := r.db.Query(query, tenantID, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	var records []models.PaymentRecord
	for rows.Next() {
		var record models.PaymentRecord
		var createdAt, updatedAt string

		err := rows.Scan(
			&record.ID,
			&record.TenantID,
			&record.Month,
			&record.TargetColdRent,
			&record.PaidColdRent,
			&record.PaidAncillary,
			&record.PaidElectricity,
			&record.ExtraPayments,
			&record.Persons,
			&record.Note,
			&record.IsLocked,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamps
		record.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		record.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Reverse the order so the oldest records come first
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}

	return records, nil
}

// GetTenantRecordsForMonthRange returns payment records for a tenant within a specified month range
func (r *PaymentRecordRepository) GetTenantRecordsForMonthRange(tenantID int64, startMonth, endMonth string) ([]models.PaymentRecord, error) {
	// Prepare the SQL statement
	query := `
		SELECT pr.id, pr.tenant_id, pr.month, pr.target_cold_rent, pr.paid_cold_rent,
			pr.paid_ancillary, pr.paid_electricity, pr.extra_payments, pr.persons,
			pr.note, pr.is_locked, pr.created_at, pr.updated_at
		FROM payment_records pr
		WHERE pr.tenant_id = ? AND pr.month >= ? AND pr.month <= ?
		ORDER BY pr.month ASC
	`

	// Execute the query
	rows, err := r.db.Query(query, tenantID, startMonth, endMonth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	var records []models.PaymentRecord
	for rows.Next() {
		var record models.PaymentRecord
		var createdAt, updatedAt string

		err := rows.Scan(
			&record.ID,
			&record.TenantID,
			&record.Month,
			&record.TargetColdRent,
			&record.PaidColdRent,
			&record.PaidAncillary,
			&record.PaidElectricity,
			&record.ExtraPayments,
			&record.Persons,
			&record.Note,
			&record.IsLocked,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamps
		record.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		record.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// GetCurrentTenantsPaymentsByHouseID gets payment records for current tenants in a specific house
func (r *PaymentRecordRepository) GetCurrentTenantsPaymentsByHouseID(houseID int64) (map[int64][]models.PaymentRecord, error) {
	// First get all current tenants for this house
	query := `
		SELECT t.id, t.first_name, t.last_name, t.move_in_date, t.apartment_id
		FROM tenants t
		WHERE t.house_id = ? AND (t.move_out_date IS NULL OR t.move_out_date > date('now'))
		ORDER BY t.last_name, t.first_name
	`

	rows, err := r.db.Query(query, houseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get current tenants IDs
	var tenantIDs []int64
	for rows.Next() {
		var tenantID int64
		var firstName, lastName, moveInDate string
		var apartmentID int64

		err := rows.Scan(&tenantID, &firstName, &lastName, &moveInDate, &apartmentID)
		if err != nil {
			return nil, err
		}

		tenantIDs = append(tenantIDs, tenantID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(tenantIDs) == 0 {
		return make(map[int64][]models.PaymentRecord), nil
	}

	// Now get payment records for these tenants
	result := make(map[int64][]models.PaymentRecord)
	for _, tenantID := range tenantIDs {
		records, err := r.GetRecentByTenantID(tenantID, 12) // Get last 12 months
		if err != nil {
			// If there's an error, just continue with the next tenant
			continue
		}
		result[tenantID] = records
	}

	return result, nil
}

// Update modifies an existing payment record in the database
func (r *PaymentRecordRepository) Update(record *models.PaymentRecord) error {
	// Validate payment record data
	if err := record.Validate(); err != nil {
		return err
	}

	// Ensure record exists
	existingRecord, err := r.GetByID(record.ID)
	if err != nil {
		return err
	}

	// Check if the record is locked
	if existingRecord.IsLocked && !record.IsLocked {
		// Allow unlocking a record
	} else if existingRecord.IsLocked && record.IsLocked {
		// If record is already locked and still locked, only allow updating the note
		record.PaidColdRent = existingRecord.PaidColdRent
		record.PaidAncillary = existingRecord.PaidAncillary
		record.PaidElectricity = existingRecord.PaidElectricity
		record.ExtraPayments = existingRecord.ExtraPayments
		record.Persons = existingRecord.Persons
	}

	// Prepare the SQL statement
	query := `
		UPDATE payment_records
		SET paid_cold_rent = ?, paid_ancillary = ?, paid_electricity = ?,
			extra_payments = ?, persons = ?, note = ?, is_locked = ?, updated_at = ?
		WHERE id = ?
	`

	// Execute the query
	now := time.Now()
	_, err = r.db.Exec(
		query,
		record.PaidColdRent,
		record.PaidAncillary,
		record.PaidElectricity,
		record.ExtraPayments,
		record.Persons,
		record.Note,
		record.IsLocked,
		now,
		record.ID,
	)
	if err != nil {
		return err
	}

	record.UpdatedAt = now

	return nil
}

// Delete removes a payment record from the database
func (r *PaymentRecordRepository) Delete(id int64) error {
	// Ensure record exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	// Prepare the SQL statement
	query := `DELETE FROM payment_records WHERE id = ?`

	// Execute the query
	_, err = r.db.Exec(query, id)
	return err
}

// BatchCreateOrUpdateRecords creates or updates multiple payment records
func (r *PaymentRecordRepository) BatchCreateOrUpdateRecords(records []models.PaymentRecord) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Try to execute all operations
	for i := range records {
		record := &records[i]

		// Check if record exists
		var id int64
		err := tx.QueryRow("SELECT id FROM payment_records WHERE tenant_id = ? AND month = ?",
			record.TenantID, record.Month).Scan(&id)

		if err != nil {
			if err == sql.ErrNoRows {
				// Check if tenant exists before creating a new record
				var tenantExists int
				err := tx.QueryRow("SELECT COUNT(*) FROM tenants WHERE id = ?", record.TenantID).Scan(&tenantExists)
				if err != nil {
					tx.Rollback()
					return err
				}
				if tenantExists == 0 {
					tx.Rollback()
					return errors.New("tenant not found")
				}

				// Create new record
				query := `
					INSERT INTO payment_records (
						tenant_id, month, target_cold_rent, paid_cold_rent, paid_ancillary,
						paid_electricity, extra_payments, persons, note, is_locked,
						created_at, updated_at
					) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				`
				now := time.Now()
				result, err := tx.Exec(
					query,
					record.TenantID,
					record.Month,
					record.TargetColdRent,
					record.PaidColdRent,
					record.PaidAncillary,
					record.PaidElectricity,
					record.ExtraPayments,
					record.Persons,
					record.Note,
					record.IsLocked,
					now,
					now,
				)
				if err != nil {
					tx.Rollback()
					return err
				}

				id, err := result.LastInsertId()
				if err != nil {
					tx.Rollback()
					return err
				}
				record.ID = id
				record.CreatedAt = now
				record.UpdatedAt = now
			} else {
				tx.Rollback()
				return err
			}
		} else {
			// Update existing record
			query := `
				UPDATE payment_records
				SET paid_cold_rent = ?, paid_ancillary = ?, paid_electricity = ?,
					extra_payments = ?, persons = ?, note = ?, is_locked = ?, updated_at = ?
				WHERE id = ?
			`

			now := time.Now()
			_, err = tx.Exec(
				query,
				record.PaidColdRent,
				record.PaidAncillary,
				record.PaidElectricity,
				record.ExtraPayments,
				record.Persons,
				record.Note,
				record.IsLocked,
				now,
				id,
			)
			if err != nil {
				tx.Rollback()
				return err
			}
			record.ID = id
			record.UpdatedAt = now
		}
	}

	return tx.Commit()
}

// GeneratePaymentRecordsForTenant creates payment records for a tenant based on their move-in date
// This is useful when a new tenant is added and we want to generate payment records automatically
func (r *PaymentRecordRepository) GeneratePaymentRecordsForTenant(tenantID int64) error {
	// First, get the tenant details
	query := `
		SELECT t.target_cold_rent, t.move_in_date, t.move_out_date, t.number_of_persons
		FROM tenants t
		WHERE t.id = ?
	`

	var targetColdRent float64
	var moveInDateStr string
	var moveOutDateStr sql.NullString
	var numberOfPersons int

	err := r.db.QueryRow(query, tenantID).Scan(
		&targetColdRent,
		&moveInDateStr,
		&moveOutDateStr,
		&numberOfPersons,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("tenant not found")
		}
		return err
	}

	// Parse move-in date
	moveInDate, err := time.Parse("2006-01-02", moveInDateStr)
	if err != nil {
		return errors.New("invalid move-in date format")
	}

	// Get the first day of the move-in month
	startMonth := time.Date(moveInDate.Year(), moveInDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Determine end month (current month or move-out month)
	var endMonth time.Time
	if moveOutDateStr.Valid {
		moveOutDate, err := time.Parse("2006-01-02", moveOutDateStr.String)
		if err != nil {
			return errors.New("invalid move-out date format")
		}
		endMonth = time.Date(moveOutDate.Year(), moveOutDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else {
		// If no move-out date, use current month
		now := time.Now()
		endMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	}

	// Generate payment records for each month
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for month := startMonth; !month.After(endMonth); month = month.AddDate(0, 1, 0) {
		monthStr := month.Format("2006-01")

		// Check if record already exists for this month
		var count int
		err := tx.QueryRow("SELECT COUNT(*) FROM payment_records WHERE tenant_id = ? AND month = ?",
			tenantID, monthStr).Scan(&count)
		if err != nil {
			tx.Rollback()
			return err
		}
		if count > 0 {
			continue // Skip if record already exists
		}

		// Insert a new payment record
		query := `
			INSERT INTO payment_records (
				tenant_id, month, target_cold_rent, paid_cold_rent, paid_ancillary,
				paid_electricity, extra_payments, persons, note, is_locked,
				created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`

		now := time.Now()
		_, err = tx.Exec(
			query,
			tenantID,
			monthStr,
			targetColdRent,
			0.0, // Default values
			0.0,
			0.0,
			0.0,
			numberOfPersons,
			"",    // Empty note
			false, // Not locked
			now,
			now,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetMonthsRange returns a slice of month strings (in YYYY-MM format) between startMonth and endMonth
func GetMonthsRange(startMonth, endMonth string) ([]string, error) {
	start, err := time.Parse("2006-01", startMonth)
	if err != nil {
		return nil, errors.New("invalid start month format")
	}

	end, err := time.Parse("2006-01", endMonth)
	if err != nil {
		return nil, errors.New("invalid end month format")
	}

	if end.Before(start) {
		return nil, errors.New("end month cannot be before start month")
	}

	var months []string
	for current := start; !current.After(end); current = current.AddDate(0, 1, 0) {
		months = append(months, current.Format("2006-01"))
	}

	return months, nil
}

// GetLast12MonthsFromToday returns a slice of the last 12 months (in YYYY-MM format)
func GetLast12MonthsFromToday() []string {
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	var months []string
	for i := 0; i < 12; i++ {
		month := currentMonth.AddDate(0, -i, 0)
		months = append(months, month.Format("2006-01"))
	}

	// Reverse the order so the oldest month comes first
	for i, j := 0, len(months)-1; i < j; i, j = i+1, j-1 {
		months[i], months[j] = months[j], months[i]
	}

	return months
}

// EnsurePaymentRecordsForActiveTenants makes sure all active tenants have payment records for the recent months
func (r *PaymentRecordRepository) EnsurePaymentRecordsForActiveTenants() error {
	// Get all active tenants
	query := `
		SELECT t.id, t.target_cold_rent, t.move_in_date, t.number_of_persons
		FROM tenants t
		WHERE t.move_out_date IS NULL OR t.move_out_date > date('now')
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Current month
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	currentMonthStr := currentMonth.Format("2006-01")

	// Process each tenant
	for rows.Next() {
		var tenantID int64
		var targetColdRent float64
		var moveInDateStr string
		var numberOfPersons int

		err := rows.Scan(
			&tenantID,
			&targetColdRent,
			&moveInDateStr,
			&numberOfPersons,
		)
		if err != nil {
			return err
		}

		// Parse move-in date
		moveInDate, err := time.Parse("2006-01-02", moveInDateStr)
		if err != nil {
			continue // Skip if invalid date
		}

		// Get the first day of the move-in month
		moveInMonth := time.Date(moveInDate.Year(), moveInDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		moveInMonthStr := moveInMonth.Format("2006-01")

		// Get existing payment records for this tenant
		existingRecords, err := r.GetByTenantID(tenantID)
		if err != nil {
			continue // Skip if error
		}

		// Create a map of existing month records
		existingMonths := make(map[string]bool)
		for _, record := range existingRecords {
			existingMonths[record.Month] = true
		}

		// Calculate how many months we need to check
		var monthsToCheck []string
		if strings.Compare(moveInMonthStr, currentMonthStr) > 0 {
			// If move-in month is in the future, skip
			continue
		} else {
			// Get months between move-in and current month
			monthRange, err := GetMonthsRange(moveInMonthStr, currentMonthStr)
			if err != nil {
				continue // Skip if error
			}
			monthsToCheck = monthRange
		}

		// Generate records for missing months
		tx, err := r.db.Begin()
		if err != nil {
			return err
		}

		for _, monthStr := range monthsToCheck {
			if !existingMonths[monthStr] {
				// Create record for this month
				query := `
					INSERT INTO payment_records (
						tenant_id, month, target_cold_rent, paid_cold_rent, paid_ancillary,
						paid_electricity, extra_payments, persons, note, is_locked,
						created_at, updated_at
					) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				`

				now := time.Now()
				_, err = tx.Exec(
					query,
					tenantID,
					monthStr,
					targetColdRent,
					0.0, // Default values
					0.0,
					0.0,
					0.0,
					numberOfPersons,
					"",    // Empty note
					false, // Not locked
					now,
					now,
				)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}
