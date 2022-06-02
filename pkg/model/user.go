package model

import (
	"strconv"
	"strings"
)

type UserData struct {
	Username      string
	State         string
	Points        int
	ConnectedIp   string
	ConnectedPort int
}

func PrintHallOfFame(users []UserData) string {
	var sb strings.Builder
	for _, user := range users {
		sb.Write([]byte(user.Username + ": " + strconv.Itoa(user.Points) + "\n"))
	}

	return sb.String()
}
