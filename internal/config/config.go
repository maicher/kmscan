package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
}

func New(path, kmscanrcExample string) (Config, error) {
	if path != "" {
		return parseConfig(path)
	}

	if dir, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		path = filepath.Join(dir, "kmscan/kmscanrc.toml")

		if fileExists(path) {
			return parseConfig(path)
		}
	}

	if dir, ok := os.LookupEnv("HOME"); ok {
		path = filepath.Join(dir, ".config/kmscan/kmscanrc.toml")

		if fileExists(path) {
			return parseConfig(path)
		}

	}

	return parseDefaultConfig(kmscanrcExample)
}

func parseDefaultConfig(kmscanrcExample string) (Config, error) {
	var c Config

	err := toml.Unmarshal([]byte(kmscanrcExample), &c)
	if err != nil {
		return c, fmt.Errorf("unable to parse default config file: %s", err)
	}

	return c, nil
}

func parseConfig(path string) (Config, error) {
	var c Config
	bytes, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return c, fmt.Errorf("invalid path to the config file")
	} else if err != nil {
		return c, fmt.Errorf("unable to read config file %s: %s", path, err)
	}

	err = toml.Unmarshal(bytes, &c)
	if err != nil {
		return c, fmt.Errorf("unable to parse default config file: %s", err)
	}

	return c, nil
}

func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
