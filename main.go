package main

import (
	"fmt"
	"log"
	"net/http"

	"csv_extractor/db"
	"csv_extractor/handlers"
)

func main() {
	http.HandleFunc("GET /healthcheck", handlers.HealthCheckHandler)
	http.HandleFunc("POST /upload", handlers.FileUploadHandler)
	http.HandleFunc("GET /categories", handlers.GetCategories)
	http.HandleFunc("POST /category", handlers.SaveCategory)
	http.HandleFunc("PUT /category", handlers.UpdateCategory)
	http.HandleFunc("DELETE /category/{id}", handlers.DisableCategory)
	http.HandleFunc("GET /expenses", handlers.GetAllExpsenses)
	http.HandleFunc("POST /expense", handlers.SaveExpense)
	http.HandleFunc("PUT /expense", handlers.UpdateExpense)
	http.HandleFunc("DELETE /expense/{id}", handlers.DisableExpense)

	err := db.Connect()

	if err != nil {
		log.Fatal("error - failed database connection: ", err)
	}

	if db.Database == nil {
		log.Fatal("error - no database connection")
	}

	defer db.Database.Close()

	fmt.Println("Server is running at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
