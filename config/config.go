package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
	"gopkg.in/yaml.v3"
)

const (
	LocalEnv    = "local"
	AppEnv      = "APP_ENV"
	DefaultPath = "./config/config."

	ClientAppID    = "CLIENT_APP_ID"
	ClientAppHash  = "CLIENT_APP_HASH"
	ClientPhone    = "CLIENT_PHONE"
	ClientPassword = "CLIENT_PASSWORD"
)

type Config struct {
	Client     Client
	GRPCServer GRPCServer `yaml:"grpc_server"`
	HTTPServer HTTPServer `yaml:"http_server"`
}
type Client struct {
	AppID    int    `yaml:"app_id"`
	AppHash  string `yaml:"app_hash"`
	Phone    string `yaml:"phone"`
	Password string `yaml:"password"`
}

type GRPCServer struct {
	Network string `yaml:"network"`
	Address string `yaml:"address"`
}

type HTTPServer struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func NewConfig() *Config {
	return &Config{}
}
func LoadConfig() (*Config, error) {

	config := NewConfig()

	env := os.Getenv(AppEnv)
	if env == "" {
		env = LocalEnv
	}
	configPath := fmt.Sprintf("%s%s.yaml", DefaultPath, env)

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config.yaml: %v", err)
	}
	err = yaml.Unmarshal(configFile, &config)

	if err != nil {
		return nil, fmt.Errorf("error parsing config.yaml: %v", err)
	}

	config.Client.AppID, err = getIntEnv(ClientAppID, config.Client.AppID)
	if err != nil {
		return nil, fmt.Errorf("error getting env: %v", err)
	}
	config.Client.AppHash = getEnv(ClientAppHash, config.Client.AppHash)
	config.Client.Phone = getEnv(ClientPhone, config.Client.Phone)
	config.Client.Password = getEnv(ClientPassword, config.Client.Password)

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int) (int, error) {

	if value, exists := os.LookupEnv(key); exists {
		value = strings.Trim(value, "\"")
		return strconv.Atoi(value)
	}
	return fallback, nil
}
