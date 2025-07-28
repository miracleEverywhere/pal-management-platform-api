package utils

import (
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"os"
	"pal-management-platform-api/store"
	"runtime"
	"strings"
)

func BindFlags() {
	flag.IntVar(&BindPort, "l", 80, "监听端口，如： -l 8080 (Listening Port, e.g. -l 8080)")
	flag.StringVar(&ConfDir, "s", "./", "数据库文件目录，如： -s ./conf (Database Directory, e.g. -s ./conf)")
	flag.BoolVar(&VersionShow, "v", false, "查看版本，如： -v (Check version, e.g. -v)")
	flag.Parse()
}

// EnsureDirExists 检查目录是否存在，如果不存在则创建
func EnsureDirExists(dirPath string) error {
	if strings.HasPrefix(dirPath, "~") {
		dirPath = strings.Replace(dirPath, "~", HomeDir, 1)
	}
	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("无法创建目录: %w", err)
		}
	} else if err != nil {
		// 其他错误
		return fmt.Errorf("检查目录时出错: %w", err)
	}

	return nil
}

func SetGlobalVariables() {
	config, err := store.ReadConfig()
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
	var err error
	// dst config
	err = EnsureDirExists(LogsPath)
	if err != nil {
		Logger.Error("目录检查未通过，创建目录失败", "path", LogsPath)
		panic("目录检查未通过")
	}
}
