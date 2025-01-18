package db

import (
	"context"
	"csv_extractor/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

func GetAllCategories(db *sql.DB, isActive bool) ([]models.Category, error) {
	rows, err := db.Query("SELECT * FROM categories WHERE is_active=$1", isActive)

	if err != nil {
		log.Fatal("Error - database querying: ", err)
		return nil, err
	}

	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var c models.Category

		err := rows.Scan(&c.Id, &c.Name, &c.Active)

		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	return categories, nil
}

func GetCategoryById(db *sql.DB, categoryId int) (*models.Category, error) {
	query := "SELECT id, name FROM categories WHERE id = $1"

	var c models.Category

	err := db.QueryRow(query, categoryId).Scan(&c.Id, &c.Name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("category not found")
		}

		return nil, err
	}

	return &c, nil
}

func CategoryExists(tx *sql.Tx, c *models.Category) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM categories WHERE name ILIKE $1 OR id = $2)"

	err := tx.QueryRow(query, c.Name, c.Id).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("error - failed to verify if category exists: %v", err)
	}

	return exists, err
}

func SaveCategory(db *sql.DB, c *models.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("error - failed to start transaction: %w", err)
	}

	defer tx.Rollback()

	exists, err := CategoryExists(tx, c)

	if err != nil {
		return fmt.Errorf("error - failed check if expense exists: %w", err)
	}

	if exists {
		return errors.New("error - expense already exists")
	}

	query := "INSERT INTO categories (name) VALUES ($1) RETURNING id"

	var id int

	err = tx.QueryRowContext(ctx, query, c.Name).Scan(&id)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error - failed to commit transaction: %w", err)
	}

	c.Id = id

	return nil
}

func UpdateCategory(db *sql.DB, c *models.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		log.Fatal("error - failed to start transaction:", err)
		return err
	}

	defer tx.Rollback()

	exists, err := CategoryExists(tx, c)

	if err != nil {
		log.Fatal("error - failed check if expense exists:", err)
		return err
	}

	if !exists {
		return errors.New("error - expense doesn't exists")
	}

	query := "UPDATE categories SET name = $1, is_active = $2 WHERE id = $3"

	res, err := tx.ExecContext(ctx, query, c.Name, c.Active, c.Id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return errors.New("error - failed row verification")
	}

	if rowsAffected == 0 {
		return errors.New("error - No lines changed")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error - failed to commit transaction: %w", err)
	}

	return nil
}
