package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var db *sql.DB

// Конфигурация базы данных
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "GoBD"
)

func main() {
	// Подключение к базе данных
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close()

	// Создание таблицы пользователей, если её нет
	createUserTable()

	// Роуты для HTTP API
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	// Запуск HTTP сервера на порту 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createUserTable() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE,
		password VARCHAR(100))`)
	if err != nil {
		log.Fatal("Ошибка создания таблицы пользователей:", err)
	}
}

// Регистрация нового пользователя
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec(
		"INSERT INTO users (username, password) VALUES ($1, $2)",
		user.Username, user.Password)
	if err != nil {
		http.Error(w, "Пользователь с таким логином уже существует", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Аутентификация пользователя
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	row := db.QueryRow(
		"SELECT id FROM users WHERE username = $1 AND password = $2",
		user.Username, user.Password)

	var userID int
	err = row.Scan(&userID)
	if err != nil {
		row := db.QueryRow(
			"SELECT COUNT(*) FROM users WHERE username = $1",
			user.Username)
		err = row.Scan()
		if err != nil {
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Ошибка аутентификации", http.StatusInternalServerError)
		return
	}

	// Отправка ID пользователя в ответе
	response := map[string]int{"userID": userID}
	json.NewEncoder(w).Encode(response)
}
