package store

import "sync"

var ConfigMutex sync.Mutex

type User struct {
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Role     string   `json:"role"`
	Avatar   string   `json:"avatar"`
	Password string   `json:"password"`
	Disabled bool     `json:"disabled"`
	Clusters []string `json:"clusters"`
}

type Config struct {
	Users      []User `json:"users"`
	JwtSecret  string `json:"jwtSecret"`
	Registered bool   `json:"registered"`
}
