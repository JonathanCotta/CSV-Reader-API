package handlers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"csv_extractor/db"
	"csv_extractor/models"
	"csv_extractor/utils"
)

func formatString(s string) string {
	var ns = s

	if strings.Contains(s, " - Parcela") {
		idx := strings.Index(s, " - Parcela")

		ns = s[:idx]
	}

	if strings.Contains(ns, " - NuPay") {
		idx := strings.Index(ns, " - NuPay")

		ns = ns[:idx]
	}

	return ns
}

func GetCsvExpenses(file multipart.File) (map[string]models.Expense, error) {
	// create csv reader
	reader := csv.NewReader(file)

	// removing header
	_, _ = reader.Read()

	var expenses = make(map[string]models.Expense)

	for {
		row, err := reader.Read()

		if err != nil {
			if err.Error() == "EOF" {
				break
			}

			fmt.Println("Line reading error", err)

			return expenses, errors.New("line reading error")
		}

		title := formatString(row[1])

		if title == "Pagamento recebido" {
			continue
		}

		value, err := strconv.ParseFloat(row[2], 64)

		if err != nil {
			fmt.Println("Value conversion error:", err)
			return expenses, errors.New("value conversion error")
		}

		if expense, ok := expenses[title]; ok {
			expense.Value += value
			expenses[title] = expense
		} else {
			e := models.Expense{
				Title: title,
				Value: value,
			}

			expenses[title] = e
		}
	}

	return expenses, nil
}

func GetExpensesCategories(es map[string]models.Expense) error {
	var hasNewExpense bool
	for t, e := range es {
		eData, err := db.GetExpenseByTitle(db.Database, e.Title)

		if err != nil {
			return err
		}

		if eData == nil {
			hasNewExpense = true
			continue
		}

		e.Id = eData.Id
		e.Active = eData.Active
		e.CategoryId = eData.CategoryId
		e.Category = eData.Category

		es[t] = e
	}

	if hasNewExpense {
		err := db.SaveExpensesBatch(db.Database, es)

		if err != nil {
			return err
		}
	}

	return nil
}

func CsvUploadHandler(w http.ResponseWriter, r *http.Request) {
	// get form data
	err := r.ParseMultipartForm(32 << 20)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, h, err := r.FormFile("file")

	if err != nil {
		utils.ErrorResponse(w, "Error: File upload", http.StatusBadRequest)
		return
	}

	if ct := h.Header.Get("Content-Type"); ct != "text/csv" {
		utils.ErrorResponse(w, "Error: File isn't a csv", http.StatusUnsupportedMediaType)
		return
	}

	defer file.Close()

	// extract expenses from csv file
	expenses, err := GetCsvExpenses(file)

	if err != nil {
		utils.ErrorResponse(w, "Error: reading file", http.StatusInternalServerError)
		return
	}

	// get/insert expenses from database
	err = GetExpensesCategories(expenses)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.DataResponse(w, "success", expenses)
}
