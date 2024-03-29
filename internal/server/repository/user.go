package repository

import (
	"database/sql"
	"ep2/pkg/model"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"time"
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
		log.Printf("A conexão ao banco falhou: %s\n", err.Error())
		panic(err)
	}

	return &UserRepository{
		db: dbPool,
	}
}

func (r *UserRepository) GetUser(name string) (*model.UserData, error) {
	getUserQuery := "SELECT name, points, state, ip FROM players " +
		"WHERE name = $1;"
	row := r.db.QueryRow(getUserQuery, name)

	var u model.UserData
	err := row.Scan(&u.Username, &u.Points, &u.State, &u.Address)
	if err != nil {
		log.Printf("Algo de errado ocorreu ao escanear um usuário do banco: %s", err.Error())
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Create(name string, password string) error {
	createUserQuery := "INSERT INTO players (id, name, password, state, points) " +
		"VALUES ($1, $2, $3, $4, $5);"
	statement, err := r.db.Prepare(createUserQuery)
	if err != nil {
		log.Printf("Algo de errado aconteceu ao inserir um novo usuário no banco: %s\n", err.Error())
		return err
	}

	_, err = statement.Exec(uuid.New().String(), name, password, "offline", 0)
	if err != nil {
		log.Printf("Algo de errado aconteceu ao inserir um novo usuário no banco: %s\n", err.Error())
		return err
	}

	log.Printf("Novo usuário <%s> foi criado.\n", name)
	return nil
}

func (r *UserRepository) ChangeStatusWithoutAddress(name string, status string) error {
	loginQuery :=
		"UPDATE players " +
			" SET state   = $1 " +
			"WHERE name   = $2 "

	_, err := r.db.Exec(loginQuery, status, name)
	if err != nil {
		log.Printf("Algo de errado aconteceu ao trocar o status do usuário <%s> no banco: %s\n", name, err.Error())
		return err
	}

	log.Printf("O usuário <%s> trocou de status para <%s>.\n", name, status)
	return nil
}

func (r *UserRepository) ChangeStatus(name string, address string, status string) error {
	loginQuery :=
		"UPDATE players " +
			" SET     state = $1, " +
			"            ip = $2, " +
			"last_heartbeat = $3  " +
			"WHERE name   = $4 "

	_, err := r.db.Exec(loginQuery, status, address, time.Now().UTC(), name)
	if err != nil {
		log.Printf("Algo de errado aconteceu ao trocar o status do usuário <%s> no banco: %s\n", name, err.Error())
		return err
	}

	log.Printf("O usuário <%s> trocou de status para <%s>.\n", name, status)
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
			log.Printf("Houve uma tentativa de recuperar a senha de um usuário %s inexistente.\n", name)
			return "", err
		}

		log.Printf("Algo de errado aconteceu ao tentar recuperar a senha atual do usuário <%s>: %s\n", name, err.Error())
		return "", err
	}

	return currentPassword, nil
}

func (r *UserRepository) HallOfFame() ([]model.UserData, error) {
	hallOfFameQuery := "SELECT name, points FROM players " +
		"ORDER BY points DESC;"
	rows, err := r.db.Query(hallOfFameQuery)
	if err != nil {
		log.Printf("Algo de errado ocorreu ao recuperar o Hall of fame do banco: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var users []model.UserData
	for rows.Next() {
		var u model.UserData
		err := rows.Scan(&u.Username, &u.Points)
		if err != nil {
			log.Printf("Algo de errado ocorreu ao escanear as linhas do Hall of Fame: %s", err.Error())
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) GetOnlineUsers() ([]model.UserData, error) {
	listOnlineQuery := "SELECT name, ip, state, last_heartbeat FROM players " +
		"WHERE state in ('online-available', 'online-playing') " +
		"ORDER BY state;"
	rows, err := r.db.Query(listOnlineQuery)
	if err != nil {
		log.Printf("Algo de errado ocorreu ao recuperar os usuários online do banco: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var users []model.UserData
	for rows.Next() {
		var u model.UserData
		err := rows.Scan(&u.Username, &u.Address, &u.State, &u.LastHeartbeat)
		if err != nil {
			log.Printf("Algo de errado ocorreu ao escanear as linhas dos usuários conectados: %s", err.Error())
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
		log.Printf("Algo de errado aconteceu ao atualizar a senha de um usuário no banco: %s\n", err.Error())
		return err
	}

	log.Printf("A senha do usuário <%s> foi atualizada.\n", name)
	return nil
}

func (r *UserRepository) Play(name1 string, name2 string, status string) error {
	playQuery := "UPDATE players " +
		"SET state = $1 " +
		"WHERE name in ($2, $3)"
	_, err := r.db.Exec(playQuery, status, name1, name2)
	if err != nil {
		log.Printf("Algo de errado aconteceu ao atualizar o status dos jogadores <%s> e <%s>: %s\n",
			name1, name2, err.Error())
		return err
	}

	return nil
}

func (r *UserRepository) UpdatePoints(name string, points int) error {
	updatePointsQuery := "UPDATE players " +
		"SET points = points + $1 " +
		"WHERE name = $2"
	_, err := r.db.Exec(updatePointsQuery, points, name)
	if err != nil {
		log.Printf("Algo de errado aconteceu ao atualizar os pontos de um usuário no banco: %s\n", err.Error())
		return err
	}

	log.Printf("Os pontos do usuário <%s> foram atualizados.\n", name)
	return nil
}

func (r *UserRepository) UpdateHeartbeats(name string, ip string) error {
	changePasswordQuery := "UPDATE players " +
		"SET ip = $1," +
		"last_heartbeat = $2 " +
		"WHERE name = $3"
	_, err := r.db.Exec(changePasswordQuery, ip, time.Now().UTC(), name)
	if err != nil {
		log.Printf("Algo de errado aconteceu ao atualizar o heartbeat de um usuário no banco: %s\n", err.Error())
		return err
	}

	log.Printf("Heartbeat recebido de <%s>!\n", name)
	return nil
}
