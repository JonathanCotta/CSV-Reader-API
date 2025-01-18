package handlers

import (
	"csv_extractor/db"
	"csv_extractor/models"
	"csv_extractor/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

/*
	TODO:
	- upload csv and return expenses with categorys
*/

func SaveExpense(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var exp models.Expense

	err := dec.Decode(&exp)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.SaveExpense(db.Database, &exp)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Successiful request")
}

func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var exp models.Expense

	err := dec.Decode(&exp)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateExpense(db.Database, &exp)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Successiful request")
}

func GetAllExpsenses(w http.ResponseWriter, r *http.Request) {
	activeStr := r.URL.Query().Get("active")

	b, err := strconv.ParseBool(activeStr)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := db.GetAllExpenses(db.Database, b)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.DataResponse(w, "Successiful request", c)
}

func GetExpense(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	idInt, _ := strconv.Atoi(id)

	exp, err := db.GetExpenseById(db.Database, idInt)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.DataResponse(w, "Successiful request", exp)
}

func DisableExpense(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	idInt, _ := strconv.Atoi(id)

	exp, err := db.GetExpenseById(db.Database, idInt)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exp.Active = false

	db.UpdateExpense(db.Database, exp)

	utils.SuccessResponse(w, "Successiful request")
}
