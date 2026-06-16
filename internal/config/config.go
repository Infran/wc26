package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

const AppName = "wc26"

type Config struct {
	API struct {
		BaseURL string `mapstructure:"base_url"`
		Timeout int    `mapstructure:"timeout"`
	} `mapstructure:"api"`
	Auth struct {
		Token  string `mapstructure:"token"`
		Email  string `mapstructure:"email"`
		Expiry string `mapstructure:"expiry"`
	} `mapstructure:"auth"`
	Output string `mapstructure:"output"`
}

func ConfigDir() (string, error) {
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	} else {
		xdg := os.Getenv("XDG_CONFIG_HOME")
		if xdg != "" {
			dir = xdg
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("finding home dir: %w", err)
			}
			dir = filepath.Join(home, ".config")
		}
	}
	return filepath.Join(dir, AppName), nil
}

func ConfigFilePath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func InitConfig() error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	viper.SetDefault("api.base_url", "https://worldcup26.ir")
	viper.SetDefault("api.timeout", 30)
	viper.SetDefault("output", "table")

	viper.SetEnvPrefix("WC26")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cfgPath, _ := ConfigFilePath()
			if err := viper.SafeWriteConfigAs(cfgPath); err != nil {
				return fmt.Errorf("writing default config: %w", err)
			}
			fmt.Fprintf(os.Stderr, "Created default config at: %s\n", cfgPath)
		} else {
			return fmt.Errorf("reading config: %w", err)
		}
	}

	return nil
}

func Set(key, value string) error {
	viper.Set(key, value)
	cfgPath, err := ConfigFilePath()
	if err != nil {
		return err
	}
	return viper.WriteConfigAs(cfgPath)
}

func Show() (string, error) {
	cfgPath, err := ConfigFilePath()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return "", fmt.Errorf("reading config: %w", err)
	}
	return string(data), nil
}
