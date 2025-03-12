package repository

import (
	"database/sql"
	"errors"
	"property-management/internal/models"
	"time"
)

// TenantRepository handles all database interactions for tenants
type TenantRepository struct {
	db *sql.DB
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

// Create adds a new tenant to the database
func (r *TenantRepository) Create(tenant *models.Tenant) error {
	// Validate tenant data
	if err := tenant.Validate(); err != nil {
		return err
	}

	// Verify that the apartment belongs to the specified house
	query := "SELECT house_id FROM apartments WHERE id = ?"
	var apartmentHouseID int64
	err := r.db.QueryRow(query, tenant.ApartmentID).Scan(&apartmentHouseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("apartment not found")
		}
		return err
	}

	if apartmentHouseID != tenant.HouseID {
		return errors.New("apartment does not belong to the specified house")
	}

	// Prepare the SQL statement
	query = `
		INSERT INTO tenants (
			first_name, last_name, move_in_date, move_out_date, deposit, 
			email, number_of_persons, target_cold_rent, target_ancillary_payment, 
			target_electricity_payment, greeting, house_id, apartment_id, 
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Format dates for SQLite
	moveInDateStr := tenant.MoveInDate.Format("2006-01-02")
	var moveOutDateStr interface{}
	if tenant.MoveOutDate != nil {
		moveOutDateStr = tenant.MoveOutDate.Format("2006-01-02")
	} else {
		moveOutDateStr = nil
	}

	// Execute the query
	now := time.Now()
	result, err := r.db.Exec(
		query,
		tenant.FirstName,
		tenant.LastName,
		moveInDateStr,
		moveOutDateStr,
		tenant.Deposit,
		tenant.Email,
		tenant.NumberOfPersons,
		tenant.TargetColdRent,
		tenant.TargetAncillaryPayment,
		tenant.TargetElectricityPayment,
		tenant.Greeting,
		tenant.HouseID,
		tenant.ApartmentID,
		now,
		now,
	)
	if err != nil {
		return err
	}

	// Get the inserted ID and update the tenant object
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	tenant.ID = id
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	return nil
}

// GetAll returns all tenants from the database
func (r *TenantRepository) GetAll() ([]models.Tenant, error) {
	// Prepare the SQL statement
	query := `
		SELECT t.id, t.first_name, t.last_name, t.move_in_date, t.move_out_date, t.deposit,
			t.email, t.number_of_persons, t.target_cold_rent, t.target_ancillary_payment,
			t.target_electricity_payment, t.greeting, t.house_id, t.apartment_id,
			t.created_at, t.updated_at,
			h.id, h.name, h.street, h.number, h.country, h.zip_code, h.city, h.created_at, h.updated_at,
			a.id, a.name, a.size, a.created_at, a.updated_at
		FROM tenants t
		JOIN houses h ON t.house_id = h.id
		JOIN apartments a ON t.apartment_id = a.id
		ORDER BY t.last_name, t.first_name
	`

	// Execute the query
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	var tenants []models.Tenant
	for rows.Next() {
		var tenant models.Tenant
		var house models.House
		var apartment models.Apartment
		var moveInDateStr, createdAtStr, updatedAtStr string
		var moveOutDateStr sql.NullString
		var houseCreatedAtStr, houseUpdatedAtStr string
		var apartmentCreatedAtStr, apartmentUpdatedAtStr string

		err := rows.Scan(
			&tenant.ID,
			&tenant.FirstName,
			&tenant.LastName,
			&moveInDateStr,
			&moveOutDateStr,
			&tenant.Deposit,
			&tenant.Email,
			&tenant.NumberOfPersons,
			&tenant.TargetColdRent,
			&tenant.TargetAncillaryPayment,
			&tenant.TargetElectricityPayment,
			&tenant.Greeting,
			&tenant.HouseID,
			&tenant.ApartmentID,
			&createdAtStr,
			&updatedAtStr,
			&house.ID,
			&house.Name,
			&house.Street,
			&house.Number,
			&house.Country,
			&house.ZipCode,
			&house.City,
			&houseCreatedAtStr,
			&houseUpdatedAtStr,
			&apartment.ID,
			&apartment.Name,
			&apartment.Size,
			&apartmentCreatedAtStr,
			&apartmentUpdatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamps and dates
		tenant.MoveInDate, _ = time.Parse("2006-01-02", moveInDateStr)
		if moveOutDateStr.Valid {
			moveOutDate, _ := time.Parse("2006-01-02", moveOutDateStr.String)
			tenant.MoveOutDate = &moveOutDate
		}
		tenant.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		tenant.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

		house.CreatedAt, _ = time.Parse(time.RFC3339, houseCreatedAtStr)
		house.UpdatedAt, _ = time.Parse(time.RFC3339, houseUpdatedAtStr)

		apartment.CreatedAt, _ = time.Parse(time.RFC3339, apartmentCreatedAtStr)
		apartment.UpdatedAt, _ = time.Parse(time.RFC3339, apartmentUpdatedAtStr)
		apartment.HouseID = house.ID

		tenant.House = &house
		tenant.Apartment = &apartment

		tenants = append(tenants, tenant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tenants, nil
}

// GetByID returns a tenant with the specified ID
func (r *TenantRepository) GetByID(id int64) (*models.Tenant, error) {
	// Prepare the SQL statement
	query := `
		SELECT t.id, t.first_name, t.last_name, t.move_in_date, t.move_out_date, t.deposit,
			t.email, t.number_of_persons, t.target_cold_rent, t.target_ancillary_payment,
			t.target_electricity_payment, t.greeting, t.house_id, t.apartment_id,
			t.created_at, t.updated_at,
			h.id, h.name, h.street, h.number, h.country, h.zip_code, h.city, h.created_at, h.updated_at,
			a.id, a.name, a.size, a.created_at, a.updated_at
		FROM tenants t
		JOIN houses h ON t.house_id = h.id
		JOIN apartments a ON t.apartment_id = a.id
		WHERE t.id = ?
	`

	// Execute the query
	var tenant models.Tenant
	var house models.House
	var apartment models.Apartment
	var moveInDateStr, createdAtStr, updatedAtStr string
	var moveOutDateStr sql.NullString
	var houseCreatedAtStr, houseUpdatedAtStr string
	var apartmentCreatedAtStr, apartmentUpdatedAtStr string

	err := r.db.QueryRow(query, id).Scan(
		&tenant.ID,
		&tenant.FirstName,
		&tenant.LastName,
		&moveInDateStr,
		&moveOutDateStr,
		&tenant.Deposit,
		&tenant.Email,
		&tenant.NumberOfPersons,
		&tenant.TargetColdRent,
		&tenant.TargetAncillaryPayment,
		&tenant.TargetElectricityPayment,
		&tenant.Greeting,
		&tenant.HouseID,
		&tenant.ApartmentID,
		&createdAtStr,
		&updatedAtStr,
		&house.ID,
		&house.Name,
		&house.Street,
		&house.Number,
		&house.Country,
		&house.ZipCode,
		&house.City,
		&houseCreatedAtStr,
		&houseUpdatedAtStr,
		&apartment.ID,
		&apartment.Name,
		&apartment.Size,
		&apartmentCreatedAtStr,
		&apartmentUpdatedAtStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}

	// Parse timestamps and dates
	tenant.MoveInDate, _ = time.Parse("2006-01-02", moveInDateStr)
	if moveOutDateStr.Valid {
		moveOutDate, _ := time.Parse("2006-01-02", moveOutDateStr.String)
		tenant.MoveOutDate = &moveOutDate
	}
	tenant.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	tenant.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

	house.CreatedAt, _ = time.Parse(time.RFC3339, houseCreatedAtStr)
	house.UpdatedAt, _ = time.Parse(time.RFC3339, houseUpdatedAtStr)

	apartment.CreatedAt, _ = time.Parse(time.RFC3339, apartmentCreatedAtStr)
	apartment.UpdatedAt, _ = time.Parse(time.RFC3339, apartmentUpdatedAtStr)
	apartment.HouseID = house.ID

	tenant.House = &house
	tenant.Apartment = &apartment

	return &tenant, nil
}

// GetByHouseID returns all tenants for a specific house
func (r *TenantRepository) GetByHouseID(houseID int64) ([]models.Tenant, error) {
	// Prepare the SQL statement
	query := `
		SELECT t.id, t.first_name, t.last_name, t.move_in_date, t.move_out_date, t.deposit,
			t.email, t.number_of_persons, t.target_cold_rent, t.target_ancillary_payment,
			t.target_electricity_payment, t.greeting, t.house_id, t.apartment_id,
			t.created_at, t.updated_at,
			h.id, h.name, h.street, h.number, h.country, h.zip_code, h.city, h.created_at, h.updated_at,
			a.id, a.name, a.size, a.created_at, a.updated_at
		FROM tenants t
		JOIN houses h ON t.house_id = h.id
		JOIN apartments a ON t.apartment_id = a.id
		WHERE t.house_id = ?
		ORDER BY t.last_name, t.first_name
	`

	// Execute the query
	rows, err := r.db.Query(query, houseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTenants(rows)
}

// GetByApartmentID returns all tenants for a specific apartment
func (r *TenantRepository) GetByApartmentID(apartmentID int64) ([]models.Tenant, error) {
	// Prepare the SQL statement
	query := `
		SELECT t.id, t.first_name, t.last_name, t.move_in_date, t.move_out_date, t.deposit,
			t.email, t.number_of_persons, t.target_cold_rent, t.target_ancillary_payment,
			t.target_electricity_payment, t.greeting, t.house_id, t.apartment_id,
			t.created_at, t.updated_at,
			h.id, h.name, h.street, h.number, h.country, h.zip_code, h.city, h.created_at, h.updated_at,
			a.id, a.name, a.size, a.created_at, a.updated_at
		FROM tenants t
		JOIN houses h ON t.house_id = h.id
		JOIN apartments a ON t.apartment_id = a.id
		WHERE t.apartment_id = ?
		ORDER BY t.last_name, t.first_name
	`

	// Execute the query
	rows, err := r.db.Query(query, apartmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTenants(rows)
}

// Update modifies an existing tenant in the database
func (r *TenantRepository) Update(tenant *models.Tenant) error {
	// Validate tenant data
	if err := tenant.Validate(); err != nil {
		return err
	}

	// Ensure tenant exists
	_, err := r.GetByID(tenant.ID)
	if err != nil {
		return err
	}

	// Verify that the apartment belongs to the specified house
	query := "SELECT house_id FROM apartments WHERE id = ?"
	var apartmentHouseID int64
	err = r.db.QueryRow(query, tenant.ApartmentID).Scan(&apartmentHouseID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("apartment not found")
		}
		return err
	}

	if apartmentHouseID != tenant.HouseID {
		return errors.New("apartment does not belong to the specified house")
	}

	// Prepare the SQL statement
	query = `
		UPDATE tenants
		SET first_name = ?, last_name = ?, move_in_date = ?, move_out_date = ?, deposit = ?,
			email = ?, number_of_persons = ?, target_cold_rent = ?, target_ancillary_payment = ?,
			target_electricity_payment = ?, greeting = ?, house_id = ?, apartment_id = ?,
			updated_at = ?
		WHERE id = ?
	`

	// Format dates for SQLite
	moveInDateStr := tenant.MoveInDate.Format("2006-01-02")
	var moveOutDateStr interface{}
	if tenant.MoveOutDate != nil {
		moveOutDateStr = tenant.MoveOutDate.Format("2006-01-02")
	} else {
		moveOutDateStr = nil
	}

	// Execute the query
	now := time.Now()
	_, err = r.db.Exec(
		query,
		tenant.FirstName,
		tenant.LastName,
		moveInDateStr,
		moveOutDateStr,
		tenant.Deposit,
		tenant.Email,
		tenant.NumberOfPersons,
		tenant.TargetColdRent,
		tenant.TargetAncillaryPayment,
		tenant.TargetElectricityPayment,
		tenant.Greeting,
		tenant.HouseID,
		tenant.ApartmentID,
		now,
		tenant.ID,
	)
	if err != nil {
		return err
	}

	tenant.UpdatedAt = now

	return nil
}

// Delete removes a tenant from the database
func (r *TenantRepository) Delete(id int64) error {
	// Ensure tenant exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	// Prepare the SQL statement
	query := `DELETE FROM tenants WHERE id = ?`

	// Execute the query
	_, err = r.db.Exec(query, id)
	return err
}

// scanTenants is a helper function to scan tenant rows
func (r *TenantRepository) scanTenants(rows *sql.Rows) ([]models.Tenant, error) {
	var tenants []models.Tenant
	for rows.Next() {
		var tenant models.Tenant
		var house models.House
		var apartment models.Apartment
		var moveInDateStr, createdAtStr, updatedAtStr string
		var moveOutDateStr sql.NullString
		var houseCreatedAtStr, houseUpdatedAtStr string
		var apartmentCreatedAtStr, apartmentUpdatedAtStr string

		err := rows.Scan(
			&tenant.ID,
			&tenant.FirstName,
			&tenant.LastName,
			&moveInDateStr,
			&moveOutDateStr,
			&tenant.Deposit,
			&tenant.Email,
			&tenant.NumberOfPersons,
			&tenant.TargetColdRent,
			&tenant.TargetAncillaryPayment,
			&tenant.TargetElectricityPayment,
			&tenant.Greeting,
			&tenant.HouseID,
			&tenant.ApartmentID,
			&createdAtStr,
			&updatedAtStr,
			&house.ID,
			&house.Name,
			&house.Street,
			&house.Number,
			&house.Country,
			&house.ZipCode,
			&house.City,
			&houseCreatedAtStr,
			&houseUpdatedAtStr,
			&apartment.ID,
			&apartment.Name,
			&apartment.Size,
			&apartmentCreatedAtStr,
			&apartmentUpdatedAtStr,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamps and dates
		tenant.MoveInDate, _ = time.Parse("2006-01-02", moveInDateStr)
		if moveOutDateStr.Valid {
			moveOutDate, _ := time.Parse("2006-01-02", moveOutDateStr.String)
			tenant.MoveOutDate = &moveOutDate
		}
		tenant.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		tenant.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)

		house.CreatedAt, _ = time.Parse(time.RFC3339, houseCreatedAtStr)
		house.UpdatedAt, _ = time.Parse(time.RFC3339, houseUpdatedAtStr)

		apartment.CreatedAt, _ = time.Parse(time.RFC3339, apartmentCreatedAtStr)
		apartment.UpdatedAt, _ = time.Parse(time.RFC3339, apartmentUpdatedAtStr)
		apartment.HouseID = house.ID

		tenant.House = &house
		tenant.Apartment = &apartment

		tenants = append(tenants, tenant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tenants, nil
}
