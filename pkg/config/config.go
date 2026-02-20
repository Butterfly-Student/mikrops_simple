package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	RBAC     RBACConfig     `mapstructure:"rbac"`
	Mikrotik MikrotikConfig `mapstructure:"mikrotik"`
	GenieACS GenieACSConfig `mapstructure:"genieacs"`
	WhatsApp WhatsAppConfig `mapstructure:"whatsapp"`
	Tripay   TripayConfig   `mapstructure:"tripay"`
	App      AppDetails     `mapstructure:"app"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Name         string `mapstructure:"name"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxLifetime  int    `mapstructure:"max_lifetime"`
}

type JWTConfig struct {
	Secret     string        `mapstructure:"secret"`
	Expiration time.Duration `mapstructure:"expiration"`
}

type MikrotikConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
}

type GenieACSConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type WhatsAppConfig struct {
	APIURL        string   `mapstructure:"api_url"`
	Token         string   `mapstructure:"token"`
	BaseURL       string   `mapstructure:"base_url"`
	APIKey        string   `mapstructure:"api_key"`
	DeviceID      string   `mapstructure:"device_id"`
	WebhookURL    string   `mapstructure:"webhook_url"`
	WebhookSecret string   `mapstructure:"webhook_secret"`
	AdminPhones   []string `mapstructure:"admin_phones"`
	UseMock       bool     `mapstructure:"use_mock"`
}

type TripayConfig struct {
	APIKey       string `mapstructure:"api_key"`
	PrivateKey   string `mapstructure:"private_key"`
	MerchantCode string `mapstructure:"merchant_code"`
	Mode         string `mapstructure:"mode"`
}

type RBACConfig struct {
	DefaultSuperAdmin DefaultSuperAdminConfig `mapstructure:"default_superadmin"`
	Casbin            CasbinConfig            `mapstructure:"casbin"`
	EventSystem       EventSystemConfig       `mapstructure:"event_system"`
	Audit             AuditConfig             `mapstructure:"audit"`
}

type DefaultSuperAdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Email    string `mapstructure:"email"`
}

type CasbinConfig struct {
	ModelFile        string `mapstructure:"model_file"`
	AutoSeedPolicies bool   `mapstructure:"auto_seed_policies"`
}

type EventSystemConfig struct {
	Type     string         `mapstructure:"type"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

type RabbitMQConfig struct {
	URL        string `mapstructure:"url"`
	Exchange   string `mapstructure:"exchange"`
	Queue      string `mapstructure:"queue"`
	RoutingKey string `mapstructure:"routing_key"`
	Durable    bool   `mapstructure:"durable"`
}

type AuditConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	LogLevel string `mapstructure:"log_level"`
}

type AppDetails struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	URL     string `mapstructure:"url"`
}

var AppConfig *Config

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.max_lifetime", 3600)
	viper.SetDefault("jwt.expiration", 3600*time.Second)
	viper.SetDefault("mikrotik.port", 8728)
	viper.SetDefault("tripay.mode", "production")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	AppConfig = &cfg
	return &cfg, nil
}

func GetConfig() *Config {
	return AppConfig
}
