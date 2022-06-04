package model

import (
	"time"
)

type UserData struct {
	Username      string
	State         string
	Address       string
	Points        int
	LastHeartbeat time.Time
}
