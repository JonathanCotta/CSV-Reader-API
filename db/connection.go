package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getDbconfig() (map[string]string, error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error: loading config file .env %v", err)
		return nil, err
	}

	//get db config
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return nil, fmt.Errorf("uma ou mais variáveis de ambiente não estão definidas")
	}

	config := map[string]string{
		"host":     host,
		"port":     port,
		"user":     user,
		"password": password,
		"dbname":   dbname,
	}

	return config, nil
}

var Database *sql.DB

func Connect() error {
	config, _ := getDbconfig()

	// String de conexão
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config["host"], config["port"], config["user"], config["password"], config["dbname"])

	var err error
	Database, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão com o banco: %w", err)
	}

	// Verifica a conexão
	if err := Database.Ping(); err != nil {
		return fmt.Errorf("erro ao verificar conexão com o banco: %w", err)
	}

	fmt.Println("Conexão com o banco estabelecida com sucesso!")
	return nil
}
