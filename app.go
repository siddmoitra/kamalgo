package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

type dbConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

func loadDbConfig() dbConfig {
	return dbConfig{
		Host:     os.Getenv("KAMALGO_DB_HOST"),
		Port:     os.Getenv("KAMALGO_DB_PORT"),
		User:     os.Getenv("KAMALGO_DB_USER"),
		Password: os.Getenv("KAMALGO_DB_PASSWORD"),
		Dbname:   os.Getenv("KAMALGO_DB_NAME"),
	}
}

func main() {
	http.HandleFunc("/", HelloWorld)
	http.HandleFunc("/up", Up)
	http.HandleFunc("/ping-db", PingDatabase)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		// Print error message to console if server fails to start
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	fmt.Println("Server started successfully")
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("<h1>Hello World!</h1>"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Up(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("<h1>All good and I am healthy!</h1>"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func PingDatabase(w http.ResponseWriter, r *http.Request) {
	ts, err := pingDb()
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("<h1>Database connection successful and DB timestamp is: <i>%s</i> !</h1>", ts)))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func pingDb() (string, error) {
	dbConf := loadDbConfig()
	fmt.Printf("Connecting to database at %s:%s with %s@%s\n", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.Dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return "", err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Error closing database connection: %v\n", err)
		} else {
			fmt.Println("Database connection closed successfully")
		}
	}(db)

	// Get DB timestamp
	var dbTimestamp string
	err = db.QueryRow("SELECT NOW()").Scan(&dbTimestamp)
	if err != nil {
		return "", fmt.Errorf("failed to get database timestamp: %v", err)
	}
	return dbTimestamp, nil
}
