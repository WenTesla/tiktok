package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"tiktok/go/storage"
)

// AppConfig 是全局应用配置，于包初始化时从 .env / 环境变量加载（环境变量优先于 .env）。
var AppConfig = loadConfig()

// Storage 是存储层单例（本地文件系统实现），于包初始化时按 StorageConfig 构造。
var Storage = storage.NewLocal(AppConfig.Storage.LocalDir, AppConfig.Storage.PublicBaseURL)

// Config 聚合所有可外部化的配置项。
type Config struct {
	MySQL  MySQLConfig
	Redis  RedisConfig
	JWT    JWTConfig
	OSS    OSSConfig
	Server ServerConfig
	Storage StorageConfig
}

type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Charset  string
	Loc      string
}

type RedisConfig struct {
	Addr     string
	Password string
}

type JWTConfig struct {
	Secret string
	Expire int // 过期时长（小时）
}

type OSSConfig struct {
	SecretId     string
	SecretKey    string
	CosURL       string
	CosURLAnd    string
	ReplaceSuffix string
}

type ServerConfig struct {
	Port int
}

type StorageConfig struct {
	Type          string // 存储类型：local（默认）
	LocalDir      string // 本地存储根目录（容器内切面，如 /data）
	PublicBaseURL string // 对外访问基址（如 http://localhost:8000/static）
}

// loadConfig 读取 .env（若存在）与环境变量，并以硬编码默认值兜底。
func loadConfig() Config {
	// 本地开发时把 .env 注入到环境变量；容器部署则直接通过 -e 注入，.env 不存在时忽略。
	_ = godotenv.Load()

	v := viper.New()
	v.AutomaticEnv()

	v.SetDefault("MYSQL_HOST", "127.0.0.1")
	v.SetDefault("MYSQL_PORT", 3306)
	v.SetDefault("MYSQL_USER", "root")
	v.SetDefault("MYSQL_PASSWORD", "zhang134679")
	v.SetDefault("MYSQL_DB", "tiktok")
	v.SetDefault("MYSQL_CHARSET", "utf8mb4")
	v.SetDefault("MYSQL_LOC", "Local")

	v.SetDefault("REDIS_ADDR", "127.0.0.1:6379")
	v.SetDefault("REDIS_PASSWORD", "")

	v.SetDefault("JWT_SECRET", "this_is_a_secret_key")
	v.SetDefault("JWT_EXPIRE_HOURS", 24*7)

	v.SetDefault("OSS_SECRET_ID", "xxxxxxxxxxxxxxxxxxx")
	v.SetDefault("OSS_SECRET_KEY", "xxxxxxxxxxxxxxxxxxx")
	v.SetDefault("OSS_COS_URL", "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com")
	v.SetDefault("OSS_COS_URL_AND", "https://tiktok-video-1313520634.cos.ap-beijing.myqcloud.com/")
	v.SetDefault("OSS_REPLACE_SUFFIX", ".jpg")

	v.SetDefault("SERVER_PORT", 8000)

	v.SetDefault("STORAGE_TYPE", "local")
	v.SetDefault("LOCAL_STORAGE_DIR", "/data")
	v.SetDefault("PUBLIC_BASE_URL", "http://localhost:8000/static")

	return Config{
		MySQL: MySQLConfig{
			Host:     v.GetString("MYSQL_HOST"),
			Port:     v.GetInt("MYSQL_PORT"),
			User:     v.GetString("MYSQL_USER"),
			Password: v.GetString("MYSQL_PASSWORD"),
			DBName:   v.GetString("MYSQL_DB"),
			Charset:  v.GetString("MYSQL_CHARSET"),
			Loc:      v.GetString("MYSQL_LOC"),
		},
		Redis: RedisConfig{
			Addr:     v.GetString("REDIS_ADDR"),
			Password: v.GetString("REDIS_PASSWORD"),
		},
		JWT: JWTConfig{
			Secret: v.GetString("JWT_SECRET"),
			Expire: v.GetInt("JWT_EXPIRE_HOURS"),
		},
		OSS: OSSConfig{
			SecretId:     v.GetString("OSS_SECRET_ID"),
			SecretKey:    v.GetString("OSS_SECRET_KEY"),
			CosURL:       v.GetString("OSS_COS_URL"),
			CosURLAnd:    v.GetString("OSS_COS_URL_AND"),
			ReplaceSuffix: v.GetString("OSS_REPLACE_SUFFIX"),
		},
		Server: ServerConfig{
			Port: v.GetInt("SERVER_PORT"),
		},
		Storage: StorageConfig{
			Type:          v.GetString("STORAGE_TYPE"),
			LocalDir:      v.GetString("LOCAL_STORAGE_DIR"),
			PublicBaseURL: v.GetString("PUBLIC_BASE_URL"),
		},
	}
}
