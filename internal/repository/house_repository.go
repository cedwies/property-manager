package repository

import (
	"database/sql"
	"errors"
	"time"

	"property-management/internal/models"
)

// HouseRepository handles all database interactions for houses
type HouseRepository struct {
	db *sql.DB
}

// NewHouseRepository creates a new house repository
func NewHouseRepository(db *sql.DB) *HouseRepository {
	return &HouseRepository{db: db}
}

// Create adds a new house to the database
func (r *HouseRepository) Create(house *models.House) error {
	// Validate house data
	if err := house.Validate(); err != nil {
		return err
	}

	// Prepare the SQL statement
	query := `
		INSERT INTO houses (name, street, number, country, zip_code, city, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Execute the query
	now := time.Now()
	result, err := r.db.Exec(
		query,
		house.Name,
		house.Street,
		house.Number,
		house.Country,
		house.ZipCode,
		house.City,
		now,
		now,
	)
	if err != nil {
		return err
	}

	// Get the inserted ID and update the house object
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	house.ID = id
	house.CreatedAt = now
	house.UpdatedAt = now

	return nil
}

// GetAll returns all houses from the database
func (r *HouseRepository) GetAll() ([]models.House, error) {
	// Prepare the SQL statement
	query := `
		SELECT id, name, street, number, country, zip_code, city, created_at, updated_at
		FROM houses
		ORDER BY name
	`

	// Execute the query
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	var houses []models.House
	for rows.Next() {
		var house models.House
		var createdAt, updatedAt string

		err := rows.Scan(
			&house.ID,
			&house.Name,
			&house.Street,
			&house.Number,
			&house.Country,
			&house.ZipCode,
			&house.City,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse timestamps
		house.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		house.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		houses = append(houses, house)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return houses, nil
}

// GetByID returns a house with the specified ID
func (r *HouseRepository) GetByID(id int64) (*models.House, error) {
	// Prepare the SQL statement
	query := `
		SELECT id, name, street, number, country, zip_code, city, created_at, updated_at
		FROM houses
		WHERE id = ?
	`

	// Execute the query
	var house models.House
	var createdAt, updatedAt string

	err := r.db.QueryRow(query, id).Scan(
		&house.ID,
		&house.Name,
		&house.Street,
		&house.Number,
		&house.Country,
		&house.ZipCode,
		&house.City,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("house not found")
		}
		return nil, err
	}

	// Parse timestamps
	house.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	house.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &house, nil
}

// Update modifies an existing house in the database
func (r *HouseRepository) Update(house *models.House) error {
	// Validate house data
	if err := house.Validate(); err != nil {
		return err
	}

	// Ensure house exists
	_, err := r.GetByID(house.ID)
	if err != nil {
		return err
	}

	// Prepare the SQL statement
	query := `
		UPDATE houses
		SET name = ?, street = ?, number = ?, country = ?, zip_code = ?, city = ?, updated_at = ?
		WHERE id = ?
	`

	// Execute the query
	now := time.Now()
	_, err = r.db.Exec(
		query,
		house.Name,
		house.Street,
		house.Number,
		house.Country,
		house.ZipCode,
		house.City,
		now,
		house.ID,
	)
	if err != nil {
		return err
	}

	house.UpdatedAt = now

	return nil
}

// Delete removes a house from the database
func (r *HouseRepository) Delete(id int64) error {
	// Ensure house exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	// Prepare the SQL statement
	query := `DELETE FROM houses WHERE id = ?`

	// Execute the query
	_, err = r.db.Exec(query, id)
	return err
}
