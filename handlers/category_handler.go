package handlers

import (
	"csv_extractor/db"
	"csv_extractor/models"
	"csv_extractor/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

func GetCategories(w http.ResponseWriter, r *http.Request) {
	activeStr := r.URL.Query().Get("active")

	b, err := strconv.ParseBool(activeStr)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := db.GetAllCategories(db.Database, b)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.DataResponse(w, "Successiful request", c)
}

func SaveCategory(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var cat models.Category

	err := dec.Decode(&cat)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	cat.Active = true

	err = db.SaveCategory(db.Database, &cat)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.DataResponse(w, "Successiful request", cat)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var cat models.Category

	err := dec.Decode(&cat)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateCategory(db.Database, &cat)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, "Successiful request")
}

func DisableCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	idInt, _ := strconv.Atoi(id)

	cat, err := db.GetCategoryById(db.Database, idInt)

	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cat.Active = false

	db.UpdateCategory(db.Database, cat)

	utils.SuccessResponse(w, "Successiful request")
}
