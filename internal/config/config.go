package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env      string `mapstructure:"ENV"`
	HTTPPort string `mapstructure:"HTTP_PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`
}

func Load() (Config, error) {
	_ = godotenv.Load()

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	_ = v.ReadInConfig()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// дефолтные знач
	v.SetDefault("ENV", "local")
	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("DB_HOST", "postgres")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_NAME", "subscriptions")
	v.SetDefault("DB_SSLMODE", "disable")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	_ = os.Setenv("TZ", "UTC")
	time.Local = time.UTC

	return cfg, nil
}

func (c Config) Addr() string {
	return fmt.Sprintf(":%s", c.HTTPPort)
}
