package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx как драйвер для database/sql
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	Email   string  `db:"email"`
	Balance float64 `db:"balance"`
}

func main() {
	// Можно переопределить через переменную окружения DB_DSN (см. README)
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://appuser:apppass@localhost:5432/appdb?sslmode=disable"
	}

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// Параметры пула
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверка соединения с таймаутом
	if err := pingWithTimeout(db, 5*time.Second); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	// ==== Демонстрация CRUD ====
	fmt.Println("== InsertUser (Dave)")
	if err := InsertUser(db, User{Name: "Dave", Email: "dave@example.com", Balance: 25}); err != nil {
		log.Printf("InsertUser: %v", err)
	}

	fmt.Println("== GetAllUsers")
	users, err := GetAllUsers(db)
	if err != nil {
		log.Fatalf("GetAllUsers: %v", err)
	}
	for _, u := range users {
		fmt.Printf("%+v\n", u)
	}

	fmt.Println("== GetUserByID(1)")
	u1, err := GetUserByID(db, 1)
	if err != nil {
		log.Printf("GetUserByID: %v", err)
	} else {
		fmt.Printf("User#1: %+v\n", u1)
	}

	// ==== Транзакция перевода средств ====
	fmt.Println("== TransferBalance Alice(1) -> Bob(2) amount=15")
	if err := TransferBalance(db, 1, 2, 15); err != nil {
		log.Printf("TransferBalance: %v", err)
	} else {
		fmt.Println("Transfer ok")
	}

	fmt.Println("== Final users:")
	users, _ = GetAllUsers(db)
	for _, u := range users {
		fmt.Printf("%+v\n", u)
	}
}

func pingWithTimeout(db *sqlx.DB, d time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	return db.PingContext(ctx)
}

// InsertUser — пример NamedExec
func InsertUser(db *sqlx.DB, user User) error {
	const q = `
		INSERT INTO users (name, email, balance)
		VALUES (:name, :email, :balance)
	`
	_, err := db.NamedExec(q, user)
	return err
}

// GetAllUsers — пример Select
func GetAllUsers(db *sqlx.DB) ([]User, error) {
	var users []User
	err := db.Select(&users, `SELECT id, name, email, balance FROM users ORDER BY id`)
	return users, err
}

// GetUserByID — пример Get
func GetUserByID(db *sqlx.DB, id int) (User, error) {
	var u User
	err := db.Get(&u, `SELECT id, name, email, balance FROM users WHERE id=$1`, id)
	return u, err
}


func TransferBalance(db *sqlx.DB, fromID, toID int, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be > 0")
	}
	if fromID == toID {
		return fmt.Errorf("sender and receiver must differ")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	// Если не закоммичено — откат
	defer func() { _ = tx.Rollback() }()

	// Блокируем и читаем отправителя
	var sender User
	if err := tx.GetContext(ctx, &sender,
		`SELECT id, name, email, balance FROM users WHERE id=$1 FOR UPDATE`, fromID); err != nil {
		return fmt.Errorf("sender not found: %w", err)
	}
	if sender.Balance < amount {
		return fmt.Errorf("insufficient funds: have %.2f, need %.2f", sender.Balance, amount)
	}

	// Дебет отправителя (защита от гонок в WHERE)
	res, err := tx.ExecContext(ctx, `
		UPDATE users
		SET balance = balance - $1
		WHERE id = $2 AND balance >= $1
	`, amount, fromID)
	if err != nil {
		return fmt.Errorf("debit sender: %w", err)
	}
	if aff, _ := res.RowsAffected(); aff != 1 {
		return fmt.Errorf("debit failed (concurrent update?)")
	}

	// Проверка существования и блокировка получателя
	var receiverExists bool
	if err := tx.GetContext(ctx, &receiverExists,
		`SELECT EXISTS(SELECT 1 FROM users WHERE id=$1 FOR UPDATE)`, toID); err != nil {
		return fmt.Errorf("check receiver: %w", err)
	}
	if !receiverExists {
		return fmt.Errorf("receiver not found")
	}

	// Кредит получателя
	if _, err := tx.ExecContext(ctx, `
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2
	`, amount, toID); err != nil {
		return fmt.Errorf("credit receiver: %w", err)
	}

	// Коммит
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}
