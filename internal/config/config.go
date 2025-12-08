package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 应用配置结构体
type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Log          LogConfig          `mapstructure:"log"`
	Database     DatabaseConfig     `mapstructure:"database"`
	Redis        RedisConfig        `mapstructure:"redis"`
	Milvus       MilvusConfig       `mapstructure:"milvus"`
	JWT          JWTConfig          `mapstructure:"jwt"`
	Auth         AuthConfig         `mapstructure:"auth"`
	CORS         CORSConfig         `mapstructure:"cors"`
	RateLimit    RateLimitConfig    `mapstructure:"ratelimit"`
	Email        EmailConfig        `mapstructure:"email"`
	Verification VerificationConfig `mapstructure:"verification"`
}

// ServerConfig HTTP 服务配置
type ServerConfig struct {
	Name         string        `mapstructure:"name"`
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

// DatabaseConfig PostgreSQL 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time"`
}

// DSN 返回数据库连接字符串
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// Addr 返回 Redis 地址
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// MilvusConfig Milvus 向量数据库配置
type MilvusConfig struct {
	Host           string `mapstructure:"host"`            // Milvus/Zilliz 地址
	Port           int    `mapstructure:"port"`            // 端口号
	CollectionName string `mapstructure:"collection_name"` // 集合名称
	Dimension      int    `mapstructure:"dimension"`       // 向量维度
	APIKey         string `mapstructure:"api_key"`         // Zilliz Cloud API Key
	UseCloud       bool   `mapstructure:"use_cloud"`       // 是否使用 Zilliz Cloud
}

// Addr 返回 Milvus 地址
func (c *MilvusConfig) Addr() string {
	// Zilliz Cloud 使用完整的 HTTPS URL
	if c.UseCloud {
		return c.Host // 直接返回完整的云端点 URL
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// JWTConfig JWT 认证配置
type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	AccessExpire  time.Duration `mapstructure:"access_expire"`
	RefreshExpire time.Duration `mapstructure:"refresh_expire"`
	Issuer        string        `mapstructure:"issuer"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	BcryptCost       int           `mapstructure:"bcrypt_cost"`
	MaxLoginAttempts int           `mapstructure:"max_login_attempts"`
	LockoutDuration  time.Duration `mapstructure:"lockout_duration"`
}

// CORSConfig CORS 跨域配置
type CORSConfig struct {
	AllowedOrigins []string      `mapstructure:"allowed_origins"`
	AllowedMethods []string      `mapstructure:"allowed_methods"`
	AllowedHeaders []string      `mapstructure:"allowed_headers"`
	MaxAge         time.Duration `mapstructure:"max_age"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled bool `mapstructure:"enabled"`
	Rate    int  `mapstructure:"rate"`
	Burst   int  `mapstructure:"burst"`
}

// EmailConfig 邮件服务配置
type EmailConfig struct {
	Provider     string `mapstructure:"provider"` // smtp | resend
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUser     string `mapstructure:"smtp_user"`
	SMTPPassword string `mapstructure:"smtp_password"`
	FromAddress  string `mapstructure:"from_address"`
	FromName     string `mapstructure:"from_name"`
	ResendAPIKey string `mapstructure:"resend_api_key"`
}

// VerificationConfig 验证码配置
type VerificationConfig struct {
	CodeLength  int           `mapstructure:"code_length"`
	CodeExpire  time.Duration `mapstructure:"code_expire"`
	ResendDelay time.Duration `mapstructure:"resend_delay"`
	MaxAttempts int           `mapstructure:"max_attempts"`
}

// global 全局配置实例
var global *Config

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	// 加载 .env 文件（如果存在）
	// godotenv 会将 .env 中的变量注入到进程环境变量中
	_ = godotenv.Load() // 忽略错误，.env 文件不存在时不影响

	v := viper.New()

	// 设置配置文件路径
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 设置环境变量前缀和替换规则
	v.SetEnvPrefix("LEXVERITAS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 显式绑定敏感环境变量（确保环境变量能覆盖配置文件）
	// 格式: 环境变量 LEXVERITAS_DATABASE_PASSWORD -> database.password
	bindEnvs := []struct {
		key     string
		envName string
	}{
		// 数据库
		{"database.password", "DATABASE_PASSWORD"},
		{"database.host", "DATABASE_HOST"},
		{"database.port", "DATABASE_PORT"},
		{"database.user", "DATABASE_USER"},
		{"database.dbname", "DATABASE_DBNAME"},
		// Redis
		{"redis.password", "REDIS_PASSWORD"},
		{"redis.host", "REDIS_HOST"},
		{"redis.port", "REDIS_PORT"},
		// JWT
		{"jwt.secret", "JWT_SECRET"},
		// 邮件
		{"email.provider", "EMAIL_PROVIDER"},
		{"email.smtp_host", "EMAIL_SMTP_HOST"},
		{"email.smtp_port", "EMAIL_SMTP_PORT"},
		{"email.smtp_user", "EMAIL_SMTP_USER"},
		{"email.smtp_password", "EMAIL_SMTP_PASSWORD"},
		{"email.resend_api_key", "EMAIL_RESEND_API_KEY"},
		// Milvus
		{"milvus.host", "MILVUS_HOST"},
		{"milvus.port", "MILVUS_PORT"},
		{"milvus.api_key", "MILVUS_API_KEY"},
		{"milvus.use_cloud", "MILVUS_USE_CLOUD"},
	}

	for _, e := range bindEnvs {
		_ = v.BindEnv(e.key, "LEXVERITAS_"+e.envName)
	}

	// 解析配置
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 设置全局配置
	global = &cfg

	return &cfg, nil
}

// Get 获取全局配置
func Get() *Config {
	return global
}
