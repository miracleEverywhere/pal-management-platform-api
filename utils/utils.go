package utils

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func BindFlags() {
	flag.IntVar(&BindPort, "l", 80, "监听端口，如： -l 8080 (Listening Port, e.g. -l 8080)")
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

func ReadContainCpuUsage() (uint64, error) {
	file, err := os.Open("/sys/fs/cgroup/cpu.stat")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "usage_usec") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return strconv.ParseUint(parts[1], 10, 64)
			}
		}
	}

	return 0, fmt.Errorf("未找到 usage_usec 数据")
}

func CpuUsage() (float64, error) {
	// 获取 CPU 使用率
	if InContainer {
		const samplingInterval = 100 * time.Millisecond // 0.1秒
		// 第一次采样
		usage1, err := ReadContainCpuUsage()
		if err != nil {
			return 0, err
		}
		// 等待 0.1 秒
		time.Sleep(samplingInterval)
		// 第二次采样
		usage2, err := ReadContainCpuUsage()
		if err != nil {
			return 0, err
		}
		// 计算 CPU 使用率百分比
		delta := usage2 - usage1
		intervalMicroseconds := float64(samplingInterval.Microseconds())
		return float64(delta) / intervalMicroseconds * 100, nil
	} else {
		percent, err := cpu.Percent(0, false)
		if err != nil {
			return 0, fmt.Errorf("error getting CPU percent: %w", err)
		}
		return percent[0], nil
	}
}

func ReadContainUintFromFile(path string) (uint64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	valueStr := strings.TrimSpace(string(data))
	if valueStr == "max" {
		return math.MaxUint64, nil
	}

	return strconv.ParseUint(valueStr, 10, 64)
}

func MemoryUsage() (float64, error) {
	// 获取内存信息
	if InContainer {
		// 读取内存限制
		data, err := os.ReadFile("/sys/fs/cgroup/memory.max")
		if err != nil {
			return 0, err
		}
		valueStr := strings.TrimSpace(string(data))
		if valueStr == "max" {
			// 没有内存限制
			vmStat, err := mem.VirtualMemory()
			if err != nil {
				return 0, fmt.Errorf("error getting virtual memory info: %w", err)
			}
			return vmStat.UsedPercent, nil
		} else {
			// 存在内存限制
			// 读取当前内存使用量
			memCurrent, err := ReadContainUintFromFile("/sys/fs/cgroup/memory.current")
			if err != nil {
				return 0, err
			}
			// 读取内存限制
			memMax, err := ReadContainUintFromFile("/sys/fs/cgroup/memory.max")
			if err != nil {
				return 0, err
			}
			return float64(memCurrent) / float64(memMax) * 100, nil
		}

	} else {
		vmStat, err := mem.VirtualMemory()
		if err != nil {
			return 0, fmt.Errorf("error getting virtual memory info: %w", err)
		}
		return vmStat.UsedPercent, nil
	}
}
