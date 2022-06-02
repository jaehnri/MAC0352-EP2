package repository

import (
	"database/sql"
	"ep2/pkg/model"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository() *UserRepository {
	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		"172.17.0.2", "postgres", "postgres", "postgres")
	dbPool, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Printf("A conexão ao banco falhou: %s\n", err.Error())
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
		fmt.Printf("Algo de errado aconteceu ao inserir um novo usuário no banco: %s\n", err.Error())
		return err
	}

	_, err = statement.Exec(uuid.New().String(), name, password, "offline", 0)
	if err != nil {
		fmt.Printf("Algo de errado aconteceu ao inserir um novo usuário no banco: %s\n", err.Error())
		return err
	}

	fmt.Printf("Novo usuário <%s> foi criado.\n", name)
	return nil
}

func (r *UserRepository) ChangeStatus(name string, address string, status string) error {
	loginQuery :=
		"UPDATE players " +
			" SET state   = $1, " +
			"     address = $2 " +
			"WHERE name   = $3 "

	_, err := r.db.Exec(loginQuery, status, address, name)
	if err != nil {
		fmt.Printf("Algo de errado aconteceu ao trocar o status do usuário <%s> no banco: %s\n", name, err.Error())
		return err
	}

	fmt.Printf("O usuário <%s> trocou de status para <%s>.\n", name, status)
	return nil
}

func (r *UserRepository) GetCurrentPassword(name string) (string, error) {
	var currentPassword string

	getOldPasswordQuery := "SELECT password FROM players " +
		"WHERE name = $1;"
	row := r.db.QueryRow(getOldPasswordQuery, name)
	err := row.Scan(&currentPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Houve uma tentativa de recuperar a senha de um usuário %s inexistente.\n", name)
			return "", err
		}

		fmt.Printf("Algo de errado aconteceu ao tentar recuperar a senha atual do usuário <%s>: %s\n", name, err.Error())
		return "", err
	}

	return currentPassword, nil
}

func (r *UserRepository) HallOfFame() ([]model.UserData, error) {
	hallOfFameQuery := "SELECT name, points FROM players " +
		"ORDER BY points DESC;"
	rows, err := r.db.Query(hallOfFameQuery)
	if err != nil {
		fmt.Printf("Algo de errado ocorreu ao recuperar o Hall of fame do banco: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var users []model.UserData
	for rows.Next() {
		var u model.UserData
		err := rows.Scan(&u.Username, &u.Points)
		if err != nil {
			fmt.Printf("Algo de errado ocorreu ao escanear as linhas do Hall of Fame: %s", err.Error())
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) GetOnlineUsers() ([]model.UserData, error) {
	listOnlineQuery := "SELECT name, address, status FROM players " +
		"WHERE status in ('online-available', 'online-playing') " +
		"ORDER BY status;"
	rows, err := r.db.Query(listOnlineQuery)
	if err != nil {
		fmt.Printf("Algo de errado ocorreu ao recuperar os usuários online do banco: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var users []model.UserData
	for rows.Next() {
		var u model.UserData
		err := rows.Scan(&u.Username, &u.Address, &u.State)
		if err != nil {
			fmt.Printf("Algo de errado ocorreu ao escanear as linhas dos usuários conectados: %s", err.Error())
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) ChangePassword(name string, password string) error {
	changePasswordQuery := "UPDATE players " +
		"SET password = $1 " +
		"WHERE name = $2"
	_, err := r.db.Exec(changePasswordQuery, password, name)
	if err != nil {
		fmt.Printf("Algo de errado aconteceu ao atualizar a senha de um usuário no banco: %s\n", err.Error())
		return err
	}

	fmt.Printf("A senha do usuário <%s> foi atualizada.\n", name)
	return nil
}
