package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

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

func CheckConfig() {
	_ = EnsureDirExists(ConfDir)
	_, err := os.Stat(ConfDir + "/DstMP.sdb")
	if !os.IsNotExist(err) {
		Logger.Info("执行数据库检查中，发现数据库文件")
		config, err := ReadConfig()
		if err != nil {
			Logger.Error("执行数据库检查中，打开数据库文件失败", "err", err)
			panic("数据库检查未通过")
			return
		}
		DBCache = config
		Logger.Info("数据库检查完成")
		return
	}

	Logger.Info("执行数据库检查中，初始化数据库")
	var config Config
	config.Init()

	Logger.Info("数据库初始化完成")
}

func (config Config) Init() {
	config.JwtSecret = GenerateJWTSecret()
	config.Registered = false
	err := WriteConfig(config)
	if err != nil {
		Logger.Error("写入数据库失败", "err", err)
		panic("数据库初始化失败")
	}
}

func ReadConfig() (Config, error) {
	if DBCache.JwtSecret != "" {
		return DBCache, nil
	}

	ConfigMutex.Lock()
	defer ConfigMutex.Unlock()

	content, err := os.ReadFile(ConfDir + "/DstMP.sdb")
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		return Config{}, fmt.Errorf("数据库格式异常: %w", err)
	}

	// 刷新缓存
	DBCache = config

	return config, nil
}

func WriteConfig(config Config) error {
	ConfigMutex.Lock()
	defer ConfigMutex.Unlock()

	data, err := json.MarshalIndent(config, "", "    ") // 格式化输出
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	file, err := os.OpenFile(ConfDir+"/DstMP.sdb", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			Logger.Error("关闭文件失败", "err", err)
		}
	}(file)

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	// 确保数据刷入磁盘
	if err := file.Sync(); err != nil {
		return fmt.Errorf("同步文件到磁盘失败: %w", err)
	}

	// 刷新缓存
	DBCache = config

	return nil
}
