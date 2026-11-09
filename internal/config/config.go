package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config armazena os parâmetros de inicialização e runtime da aplicação.
type Config struct {
	BaseURL          string        `mapstructure:"base_url"`
	Token            string        `mapstructure:"token"`
	Timeout          time.Duration `mapstructure:"timeout"`
	ColdStartTimeout time.Duration `mapstructure:"cold_start_timeout"`
	ConfigPath       string        `mapstructure:"-"`
}

const (
	DefaultBaseURL          = "https://full-cycle-web-api.onrender.com"
	DefaultTimeout          = 12 * time.Second
	DefaultColdStartTimeout = 45 * time.Second
)

// Load carrega configurações com prioridade: Flags/Env > Arquivo YAML > Padrões.
func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("base_url", DefaultBaseURL)
	v.SetDefault("timeout", DefaultTimeout)
	v.SetDefault("cold_start_timeout", DefaultColdStartTimeout)

	// Variáveis de ambiente: STAR_API_URL, STAR_API_TOKEN, etc.
	v.SetEnvPrefix("STAR")
	v.AutomaticEnv()
	_ = v.BindEnv("base_url", "STAR_API_URL")
	_ = v.BindEnv("token", "STAR_API_TOKEN")

	// Localização cross-platform do arquivo de configuração
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	appDir := filepath.Join(configDir, "star-tui")
	_ = os.MkdirAll(appDir, 0700)

	configFilePath := filepath.Join(appDir, "config.yaml")
	v.SetConfigFile(configFilePath)

	if err := v.ReadInConfig(); err != nil {
		// Se o arquivo não existir ainda, ignora o erro
		if !os.IsNotExist(err) && !os.IsPermission(err) {
			// apenas se for erro sintático
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("falha ao interpretar configuracoes: %w", err)
	}

	cfg.ConfigPath = configFilePath
	return &cfg, nil
}

// Save grava alterações no arquivo de configuração do usuário.
func (c *Config) Save() error {
	v := viper.New()
	v.SetConfigFile(c.ConfigPath)
	v.Set("base_url", c.BaseURL)
	v.Set("token", c.Token)
	v.Set("timeout", c.Timeout.String())
	v.Set("cold_start_timeout", c.ColdStartTimeout.String())

	return v.WriteConfigAs(c.ConfigPath)
}
