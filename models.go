package pixie

import (
	"net"
	"time"
)

type (
	Room struct {
	}

	Assignment struct {
		DjUser string
		Room   *Room
		Host   *Host
	}

	Host struct {
		rId Room
	}

	DjUser struct {
		Id       string     `json:"id"`
		Username string     `json:"username"`
		Name     string     `json:"name"`
		Email    string     `json:"email"`
		LIP      *net.IPNet `json:"last_ip"`
		IP       *net.IPNet `json:"ip"`
		LLT      time.Time  `json:"last_login_time"`
		FLT      time.Time  `json:"first_login_time"`
		Team     string     `json:"team"`
		Roles    []string   `json:"roles"`
		Enabled  bool       `json:"enabled"`
	}
)
