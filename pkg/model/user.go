package model

import (
	"strconv"
	"strings"
)

type UserData struct {
	Username string
	State    string
	Address  string
	Points   int
}

func PrintHallOfFame(users []UserData) string {
	var sb strings.Builder
	for _, user := range users {
		sb.Write([]byte(user.Username + ": " + strconv.Itoa(user.Points) + "\n"))
	}

	return sb.String()
}

func PrintOnlineUsers(users []UserData) string {
	var sb strings.Builder
	for _, user := range users {
		sb.Write([]byte(user.Username + ": " + user.Address + " - " + user.State))
	}

	return sb.String()
}
