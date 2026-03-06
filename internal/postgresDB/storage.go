package postgresDB

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func InitDB() *sql.DB {
	connStr := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("База не ответила на пинг: ", err)
	}

	fmt.Println("Бот подключен к базе данных!")
	return db
}

func CreateTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		tg_id BIGINT UNIQUE NOT NULL,
		username TEXT
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Не удалось создать/найти таблицу:", err)
	}
}

func SaveUser(db *sql.DB, tgID int64, username string) {
	query := `INSERT INTO users (tg_id, username) VALUES ($1, $2) ON CONFLICT (tg_id) DO NOTHING`
	_, err := db.Exec(query, tgID, username)
	if err != nil {
		log.Println("Ошибка сохранения пользователя:", err)
	}
}

func GetAllUsers(db *sql.DB) []int64 {
	rows, err := db.Query("SELECT tg_id FROM users")
	if err != nil {
		log.Println("Ошибка получения списка юзеров:", err)
		return nil
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		} else {
			log.Println("Ошибка добавления пользователя: ", id)
		}
	}
	return ids
}
