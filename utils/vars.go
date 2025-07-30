package utils

import "github.com/dgrijalva/jwt-go"

var (
	BindPort    int
	VersionShow bool
)

var (
	HomeDir     string
	Platform    string
	Registered  bool
	InContainer bool
)

type Claims struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`
	jwt.StandardClaims
}

type OSInfo struct {
	Architecture    string
	OS              string
	CPUModel        string
	CPUCores        int
	MemorySize      uint64
	Platform        string
	PlatformVersion string
	Uptime          uint64
}
