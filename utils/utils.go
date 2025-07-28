package utils

import (
	"flag"
	"github.com/dgrijalva/jwt-go"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func BindFlags() {
	flag.IntVar(&BindPort, "l", 80, "监听端口，如： -l 8080 (Listening Port, e.g. -l 8080)")
	flag.StringVar(&ConfDir, "s", "./", "数据库文件目录，如： -s ./conf (Database Directory, e.g. -s ./conf)")
	flag.BoolVar(&VersionShow, "v", false, "查看版本，如： -v (Check version, e.g. -v)")
	flag.Parse()
}

func SetGlobalVariables() {
	config, err := ReadConfig()
	if err != nil {
		Logger.Error("启动检查出现致命错误：获取数据库失败", "err", err)
		panic(err)
	}

	HomeDir, err = os.UserHomeDir()
	if err != nil {
		Logger.Error("无法获取用户HOME目录", "err", err)
		panic("无法获取用户HOME目录")
	}

	osInfo, err := GetOSInfo()
	if err != nil {
		Logger.Error("启动检查出现致命错误：获取系统信息失败", "err", err)
		panic(err)
	}
	Platform = osInfo.Platform

	Registered = config.Registered

	// 查看是否在容器内
	_, InContainer = os.LookupEnv("DMP_IN_CONTAINER")
}

func GetOSInfo() (*OSInfo, error) {
	architecture := runtime.GOARCH

	// 获取CPU信息
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	cpuModel := cpuInfo[0].ModelName
	cpuCount, _ := cpu.Counts(true)
	cpuCore := cpuCount

	// 获取内存信息
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	memorySize := virtualMemory.Total

	// 获取主机信息
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}
	platformVersion := hostInfo.PlatformVersion
	platform := hostInfo.Platform
	uptime := hostInfo.Uptime
	osName := hostInfo.OS
	// 返回系统信息
	return &OSInfo{
		Architecture:    architecture,
		OS:              osName,
		CPUModel:        cpuModel,
		CPUCores:        cpuCore,
		MemorySize:      memorySize,
		Platform:        platform,
		Uptime:          uptime,
		PlatformVersion: platformVersion,
	}, nil
}

func CheckDirs() {

}

func GenerateJWT(user User, jwtSecret []byte, expiration int) (string, error) {
	// 定义一个自定义的声明结构

	claims := Claims{
		Username: user.Username,
		Nickname: user.Nickname,
		Role:     user.Role,
		Avatar:   user.Avatar,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expiration) * time.Hour).Unix(), // 过期时间
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateJWTSecret() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 26
	randomString := make([]byte, length)
	for i := range randomString {
		// 从字符集中随机选择一个字符
		randomString[i] = charset[r.Intn(len(charset))]
	}

	return string(randomString)
}
