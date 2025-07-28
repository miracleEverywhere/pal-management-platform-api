package store

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"pal-management-platform-api/utils"
	"time"
)

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

func CheckConfig() {
	_ = utils.EnsureDirExists(utils.ConfDir)
	_, err := os.Stat(utils.ConfDir + "/DstMP.sdb")
	if !os.IsNotExist(err) {
		utils.Logger.Info("执行数据库检查中，发现数据库文件")
		config, err := ReadConfig()
		if err != nil {
			utils.Logger.Error("执行数据库检查中，打开数据库文件失败", "err", err)
			panic("数据库检查未通过")
			return
		}
		DBCache = config
		utils.Logger.Info("数据库检查完成")
		return
	}

	utils.Logger.Info("执行数据库检查中，初始化数据库")
	var config Config
	config.Init()

	utils.Logger.Info("数据库初始化完成")
}

func (config Config) Init() {
	config.JwtSecret = GenerateJWTSecret()
	config.Registered = false
	err := WriteConfig(config)
	if err != nil {
		utils.Logger.Error("写入数据库失败", "err", err)
		panic("数据库初始化失败")
	}
}

func ReadConfig() (Config, error) {
	if DBCache.JwtSecret != "" {
		return DBCache, nil
	}

	ConfigMutex.Lock()
	defer ConfigMutex.Unlock()

	content, err := os.ReadFile(utils.ConfDir + "/DstMP.sdb")
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
	file, err := os.OpenFile(utils.ConfDir+"/DstMP.sdb", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			utils.Logger.Error("关闭文件失败", "err", err)
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
