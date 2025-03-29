package config

// Config 应用配置项
type Config struct {
	Server ServerConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
}

// GetConfig 获取应用配置
func GetConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
		},
	}
}
