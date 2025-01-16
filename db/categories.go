package db

import (
	"csv_extractor/models"
	"database/sql"
	"errors"
	"log"
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
			return nil, errors.New("Category not found")
		}

		return nil, err
	}

	return &c, nil
}

func SaveCategory(db *sql.DB, n string) (int, error) {
	query := "INSERT INTO categories (name) VALUES ($1) RETURNING id"

	var id int

	err := db.QueryRow(query, n).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateCategory(db *sql.DB, c *models.Category) error {
	query := "UPDATE categories SET name = $1, is_active = $2 WHERE id = $3"

	res, err := db.Exec(query, c.Name, c.Active, c.Id)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return errors.New("Error - failed row verification")
	}

	if rowsAffected == 0 {
		return errors.New("Error - No lines changed")
	}

	return nil
}

func CategoryExists(db *sql.DB, c *models.Category) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM categories WHERE name ILIKE $1 OR id = $2)"

	err := db.QueryRow(query, c.Name, c.Id).Scan(&exists)

	return exists, err
}
