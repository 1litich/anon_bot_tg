package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func Database() {
	var err error
	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		fmt.Println("Ошибка: переменная окружения DATABASE_URL не задана")
		os.Exit(1)
	}

	db, err = sql.Open("pgx", connStr)

	if err != nil {
		fmt.Println("Ошибка настроек подключения:", err)
		os.Exit(1)
	}

	//defer db.Close()

	err = db.Ping()

	if err != nil {
		fmt.Println("Сервер базы данных недоступен!", err)
		os.Exit(1)
	}

	queryCreateMID := `CREATE TABLE IF NOT EXISTS senderid(
	id SERIAL PRIMARY KEY,
	message_id INTEGER,
	sender_id BIGINT
	);`

	_, err = db.Exec(queryCreateMID)

	if err != nil {
		fmt.Println("Ошибка создания таблицы senderid!", err)
		os.Exit(1)
	}

	queryCreateTarget := `CREATE TABLE IF NOT EXISTS targetid (
	id SERIAL PRIMARY KEY,
	target_id BIGINT,
	user_id BIGINT UNIQUE
	);`

	_, err = db.Exec(queryCreateTarget)

	if err != nil {
		fmt.Println("Ошибка создания таблицы targetid!", err)
		os.Exit(1)
	}

}

func Close() {
	db.Close()
}
