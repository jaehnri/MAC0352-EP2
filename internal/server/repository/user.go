package repository

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type User interface {
	Create(name string, password string)
	ChangePassword(name string, password string)
	Login(name string)
	Logout(name string)
	ListConnected()
	ListAll()
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		"172.17.0.2", "postgres", "postgres", "postgres")
	dbPool, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	return &UserRepository{
		db: dbPool,
	}
}

func (r *UserRepository) Create(name string, password string) error {
	createUserQuery := "INSERT INTO players (id, name, password, state, points) " +
		"VALUES ($1, $2, $3, $4, $5);"
	statement, err := r.db.Prepare(createUserQuery)
	if err != nil {
		fmt.Print(err)
		return err
	}

	_, err = statement.Exec(uuid.New().String(), name, password, "offline", 0)
	if err != nil {
		return err
	}

	fmt.Printf("Created user <%s>.\n", name)
	return nil
}
