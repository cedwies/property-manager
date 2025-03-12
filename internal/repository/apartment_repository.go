package repository

import (
	"database/sql"
	"errors"
	"time"

	"property-management/internal/models"
)

// ApartmentRepository handles all database interactions for apartments
type ApartmentRepository struct {
	db *sql.DB
}

// NewApartmentRepository creates a new apartment repository
func NewApartmentRepository(db *sql.DB) *ApartmentRepository {
	return &ApartmentRepository{db: db}
}

// Create adds a new apartment to the database
func (r *ApartmentRepository) Create(apartment *models.Apartment) error {
	// Validate apartment data
	if err := apartment.Validate(); err != nil {
		return err
	}

	// Prepare the SQL statement
	query := `
		INSERT INTO apartments (name, house_id, size, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	// Execute the query
	now := time.Now()
	result, err := r.db.Exec(
		query,
		apartment.Name,
		apartment.HouseID,
		apartment.Size,
		now,
		now,
	)
	if err != nil {
		return err
	}

	// Get the inserted ID and update the apartment object
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	apartment.ID = id
	apartment.CreatedAt = now
	apartment.UpdatedAt = now

	return nil
}

// GetAll returns all apartments from the database
func (r *ApartmentRepository) GetAll() ([]models.Apartment, error) {
	// Prepare the SQL statement
	query := `
		SELECT a.id, a.name, a.house_id, a.size, a.created_at, a.updated_at,
               h.id, h.name, h.street, h.number, h.country, h.zip_code, h.city, h.created_at, h.updated_at
		FROM apartments a
		JOIN houses h ON a.house_id = h.id
		ORDER BY a.name
	`

	// Execute the query
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	var apartments []models.Apartment
	for rows.Next() {
		var apartment models.Apartment
		var house models.House
		var createdAt, updatedAt, houseCreatedAt, houseUpdatedAt string

		err := rows.Scan(
			&apartment.ID,
			&apartment.Name,
			&apartment.HouseID,
			&apartment.Size,
			&createdAt,
			&updatedAt,
			&house.ID,
			&house.Name,
			&house.Street,
			&house.Number,
			&house.Country,
			&house.ZipCode,
			&house.City,
			&houseCreatedAt,
			&houseUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamps
		apartment.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		apartment.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		house.CreatedAt, _ = time.Parse(time.RFC3339, houseCreatedAt)
		house.UpdatedAt, _ = time.Parse(time.RFC3339, houseUpdatedAt)

		apartment.House = &house

		apartments = append(apartments, apartment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return apartments, nil
}

// GetByHouseID returns all apartments for a specific house
func (r *ApartmentRepository) GetByHouseID(houseID int64) ([]models.Apartment, error) {
	// Prepare the SQL statement
	query := `
		SELECT a.id, a.name, a.house_id, a.size, a.created_at, a.updated_at,
               h.id, h.name, h.street, h.number, h.country, h.zip_code, h.city, h.created_at, h.updated_at
		FROM apartments a
		JOIN houses h ON a.house_id = h.id
		WHERE a.house_id = ?
		ORDER BY a.name
	`

	// Execute the query
	rows, err := r.db.Query(query, houseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	var apartments []models.Apartment
	for rows.Next() {
		var apartment models.Apartment
		var house models.House
		var createdAt, updatedAt, houseCreatedAt, houseUpdatedAt string

		err := rows.Scan(
			&apartment.ID,
			&apartment.Name,
			&apartment.HouseID,
			&apartment.Size,
			&createdAt,
			&updatedAt,
			&house.ID,
			&house.Name,
			&house.Street,
			&house.Number,
			&house.Country,
			&house.ZipCode,
			&house.City,
			&houseCreatedAt,
			&houseUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamps
		apartment.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		apartment.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		house.CreatedAt, _ = time.Parse(time.RFC3339, houseCreatedAt)
		house.UpdatedAt, _ = time.Parse(time.RFC3339, houseUpdatedAt)

		apartment.House = &house

		apartments = append(apartments, apartment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return apartments, nil
}

// GetByID returns an apartment with the specified ID
func (r *ApartmentRepository) GetByID(id int64) (*models.Apartment, error) {
	// Prepare the SQL statement
	query := `
		SELECT a.id, a.name, a.house_id, a.size, a.created_at, a.updated_at,
               h.id, h.name, h.street, h.number, h.country, h.zip_code, h.city, h.created_at, h.updated_at
		FROM apartments a
		JOIN houses h ON a.house_id = h.id
		WHERE a.id = ?
	`

	// Execute the query
	var apartment models.Apartment
	var house models.House
	var createdAt, updatedAt, houseCreatedAt, houseUpdatedAt string

	err := r.db.QueryRow(query, id).Scan(
		&apartment.ID,
		&apartment.Name,
		&apartment.HouseID,
		&apartment.Size,
		&createdAt,
		&updatedAt,
		&house.ID,
		&house.Name,
		&house.Street,
		&house.Number,
		&house.Country,
		&house.ZipCode,
		&house.City,
		&houseCreatedAt,
		&houseUpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("apartment not found")
		}
		return nil, err
	}

	// Parse timestamps
	apartment.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	apartment.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	house.CreatedAt, _ = time.Parse(time.RFC3339, houseCreatedAt)
	house.UpdatedAt, _ = time.Parse(time.RFC3339, houseUpdatedAt)

	apartment.House = &house

	return &apartment, nil
}

// Update modifies an existing apartment in the database
func (r *ApartmentRepository) Update(apartment *models.Apartment) error {
	// Validate apartment data
	if err := apartment.Validate(); err != nil {
		return err
	}

	// Ensure apartment exists
	_, err := r.GetByID(apartment.ID)
	if err != nil {
		return err
	}

	// Prepare the SQL statement
	query := `
		UPDATE apartments
		SET name = ?, house_id = ?, size = ?, updated_at = ?
		WHERE id = ?
	`

	// Execute the query
	now := time.Now()
	_, err = r.db.Exec(
		query,
		apartment.Name,
		apartment.HouseID,
		apartment.Size,
		now,
		apartment.ID,
	)
	if err != nil {
		return err
	}

	apartment.UpdatedAt = now

	return nil
}

// Delete removes an apartment from the database
func (r *ApartmentRepository) Delete(id int64) error {
	// Ensure apartment exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	// Prepare the SQL statement
	query := `DELETE FROM apartments WHERE id = ?`

	// Execute the query
	_, err = r.db.Exec(query, id)
	return err
}
