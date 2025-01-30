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

func ExpenseExists(tx *sql.Tx, e *models.Expense) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM expenses WHERE title ILIKE $1 OR id = $2)"

	err := tx.QueryRow(query, e.Title, e.Id).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("error - failed to verify if expense exists: %v", err.Error())
	}

	return exists, err
}

func SaveExpense(db *sql.DB, expense *models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("error - failed to start transaction: %s", err.Error())
	}

	defer tx.Rollback()

	exists, err := ExpenseExists(tx, expense)

	if err != nil {
		return fmt.Errorf("error - failed check if expense exists: %s", err.Error())
	}

	if exists {
		return errors.New("error - expense already exists")
	}

	query := "INSERT INTO expenses (title, category_id) VALUES ($1, $2) RETURNING id"

	var id int

	err = tx.QueryRowContext(ctx, query, expense.Title, expense.CategoryId).Scan(&id)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error - failed to commit transaction: %s", err.Error())
	}

	expense.Id = id
	return nil
}

func GetAllExpenses(db *sql.DB, isActive bool) ([]models.Expense, error) {
	query := `SELECT e.id, e.title, c.id, c.name, e.is_active 
	FROM expenses e
	LEFT JOIN categories c ON e.category_id = c.id`

	rows, err := db.Query(query, isActive)

	if err != nil {
		log.Fatal("erro ao buscar: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	var expenses []models.Expense

	for rows.Next() {
		var e models.Expense

		err := rows.Scan(&e.Id, &e.Title, &e.Category, &e.CategoryId, &e.Active)

		if err != nil {
			return nil, err
		}

		expenses = append(expenses, e)
	}

	return expenses, nil
}

func GetExpenseById(db *sql.DB, expenseId int) (*models.Expense, error) {
	query := `SELECT e.id, e.title, c.id, c.name, e.is_active 
	FROM expenses e
	LEFT JOIN categories c ON e.category_id = c.id
	WHERE c.id = $1`

	var e models.Expense

	err := db.QueryRow(query, expenseId).Scan(&e.Id, &e.Title, &e.CategoryId, &e.Category, &e.Active)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("expense not found")
		}

		return nil, err
	}

	return &e, nil
}

func GetExpenseByTitle(db *sql.DB, t string) (*models.Expense, error) {
	query := `SELECT e.id, e.title, c.id, c.name, e.is_active 
	FROM expenses e
	LEFT JOIN categories c ON e.category_id = c.id
	WHERE e.title = $1`

	var e models.Expense

	err := db.QueryRow(query, t).Scan(&e.Id, &e.Title, &e.CategoryId, &e.Category, &e.Active)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &e, nil
}

func UpdateExpense(db *sql.DB, e *models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		log.Fatal("error - failed to start transaction:", err.Error())
		return err
	}

	defer tx.Rollback()

	exists, err := ExpenseExists(tx, e)

	if err != nil {
		log.Fatal("error - failed check if expense exists:", err.Error())
		return err
	}

	if !exists {
		return errors.New("error - expense doesn't exists")
	}

	query := "UPDATE expenses SET title = $1, is_active = $2, category_id = $3 WHERE id = $4"

	res, err := tx.ExecContext(ctx, query, e.Title, e.Active, e.CategoryId, e.Id)

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
		return fmt.Errorf("error - failed to commit transaction: %s", err.Error())
	}

	return nil
}

func SaveExpensesBatch(db *sql.DB, e map[string]models.Expense) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	tsx, err := db.BeginTx(ctx, nil)

	if err != nil {
		log.Fatal("error - failed to start transaction:", err.Error())
		return err
	}

	defaultCategory, err := GetCategoryByName(db, "Outros")

	if err != nil {
		log.Fatal("error - failed to get default category:", err.Error())
		return err
	}

	defer tsx.Rollback()

	stmt, err := tsx.PrepareContext(ctx, "INSERT INTO expenses (title, category_id) VALUES ($1, $2) RETURNING id")

	if err != nil {
		log.Fatal("error - failed to prepare statement:", err.Error())
		return err
	}

	defer stmt.Close()

	for _, expense := range e {
		if expense.Id != 0 {
			continue
		}

		var newId int

		err := stmt.QueryRowContext(ctx, expense.Title, defaultCategory.Id).Scan(&newId)

		if err != nil {
			return fmt.Errorf("error - failed to save expense %s: %v", expense.Title, err.Error())
		}

		expense.Id = newId
		expense.Category = defaultCategory.Name
		expense.CategoryId = defaultCategory.Id

		e[expense.Title] = expense
	}

	if err := tsx.Commit(); err != nil {
		return fmt.Errorf("error - failed to commit transaction: %s", err.Error())
	}

	return nil
}
