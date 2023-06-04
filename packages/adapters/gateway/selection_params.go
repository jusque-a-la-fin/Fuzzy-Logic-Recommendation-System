package gateway

import (
	"database/sql"
	"fmt"
	"strings"
	"vehicles/packages/domain/models"
	"vehicles/packages/usecase/repository"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type selectionRepository struct {
	preferencesDB *sql.DB
	vehiclesDB    *sql.DB
	rdb           *redis.Client
}

func NewSelectionRepository(preferencesDB *sql.DB, vehiclesDB *sql.DB, rdb *redis.Client) repository.SelectionRepository {
	return &selectionRepository{preferencesDB, vehiclesDB, rdb}
}

func (sr *selectionRepository) InsertPriorities(fingerprint string, priorities *[]string) (string, error) {

	_, err := sr.preferencesDB.Exec("INSERT INTO fingerprints (fingerprint) VALUES ($1)", fingerprint)
	if err != nil {
		return "", err
	}
	fmt.Println(priorities)

	var fingerprintID int
	err = sr.preferencesDB.QueryRow("SELECT id FROM fingerprints WHERE fingerprint = $1;", fingerprint).Scan(&fingerprintID)
	if err != nil {
		return "", fmt.Errorf("failed to get fingerprint ID: %v", err)
	}

	prioritiesArray := pq.Array(priorities)

	query := `
        INSERT INTO preferences (fingerprint_id, priorities)
        VALUES ($1, $2)
        RETURNING id;
    `
	stmt, err := sr.preferencesDB.Prepare(query)
	if err != nil {
		return "", fmt.Errorf("failed to prepare query: %v", err)
	}
	defer stmt.Close()

	var preferencesID string
	err = stmt.QueryRow(fingerprintID, prioritiesArray).Scan(&preferencesID)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %v", err)
	}

	return preferencesID, nil
}

func (sr *selectionRepository) UpdatePriorities(preferencesID, fingerprint string, priorities *[]string) error {

	var fingerprintID int
	err := sr.preferencesDB.QueryRow("SELECT id FROM fingerprints WHERE fingerprint = $1;", fingerprint).Scan(&fingerprintID)
	if err != nil {
		return fmt.Errorf("failed to get fingerprint ID: %v", err)
	}

	stmt, err := sr.preferencesDB.Prepare("UPDATE preferences SET priorities = $1 WHERE id = $2 AND fingerprint_id = $3")
	if err != nil {
		return err
	}
	defer stmt.Close()

	pgArray, err := pq.Array(priorities).Value()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(pgArray, preferencesID, fingerprintID)
	if err != nil {
		return err
	}
	return nil
}

func (sr *selectionRepository) SetPrice(preferencesID, fingerprint, minPrice, maxPrice, deviation string) error {

	var fingerprintID int
	err := sr.preferencesDB.QueryRow("SELECT id FROM fingerprints WHERE fingerprint = $1;", fingerprint).Scan(&fingerprintID)
	if err != nil {
		return fmt.Errorf("failed to get fingerprint ID: %v", err)
	}

	stmt, err := sr.preferencesDB.Prepare("UPDATE preferences SET min_price = $1, max_price = $2 WHERE id = $3 AND fingerprint_id = $4")
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	// Execute the prepared statement with the provided values
	if maxPrice == "" {
		maxPrice = "0"
	}
	if minPrice == "" {
		minPrice = "0"
	}

	res, err := stmt.Exec(minPrice, maxPrice, preferencesID, fingerprintID)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	// Check the number of rows affected by the statement (should be 1)
	numRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get row count: %v", err)
	}
	if numRows != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", numRows)
	}

	return nil
}

func (sr *selectionRepository) SetManufacturers(preferencesID, fingerprint string, manufacturers *[]string) error {

	var fingerprintID int
	err := sr.preferencesDB.QueryRow("SELECT id FROM fingerprints WHERE fingerprint = $1;", fingerprint).Scan(&fingerprintID)
	if err != nil {
		return fmt.Errorf("failed to get fingerprint ID: %v", err)
	}

	stmt, err := sr.preferencesDB.Prepare("UPDATE preferences SET manufacturers = $1 WHERE id = $2 AND fingerprint_id = $3")
	if err != nil {
		return err
	}
	defer stmt.Close()

	pgArray, err := pq.Array(manufacturers).Value()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(pgArray, preferencesID, fingerprintID)
	if err != nil {
		return err
	}
	return nil
}

func (sr *selectionRepository) GetSelection(preferencesID string) (*models.Selection, error) {

	sl := new(models.Selection)
	var priorities string
	var manufacturers string

	err := sr.preferencesDB.QueryRow(`SELECT priorities, min_price, max_price, manufacturers FROM preferences WHERE id = $1`, preferencesID).Scan(&priorities, &sl.MinPrice, &sl.MaxPrice, &manufacturers)
	if err != nil {

		return nil, err
	}

	if priorities != "{}" {
		sl.Priorities = strings.Split(priorities[1:len(priorities)-1], ",")
	}

	if manufacturers != "{}" {
		sl.Manufacturers = strings.Split(manufacturers[1:len(manufacturers)-1], ",")
	}
	return sl, nil
}
