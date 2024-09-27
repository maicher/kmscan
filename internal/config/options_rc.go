package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/guregu/null.v3"
)

type OptionsRC struct {
	Profiles []Profile `toml:"profile"`

	ProfileName null.String `toml:"profile-name"`
	CacheDir    null.String `toml:"cache-dir"`
	ResultDir   null.String `toml:"result-dir"`
	ForceDetect null.Bool   `toml:"force-detect"`
	Debug       null.Bool   `toml:"debug"`
	DebugDir    null.String `toml:"debug-dir"`
	NoColor     null.Bool   `toml:"no-color"`
	ImagePath   null.String `toml:"image-path"`
	APIURL      null.String `toml:"api-url"`
	APIKey      null.String `toml:"api-key"`
}

func NewOptionsRC(kmscanrcExample string) (*OptionsRC, error) {
	var path string

	if dir, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		path = filepath.Join(dir, "kmscan/kmscanrc.toml")

		if fileExists(path) {
			return parseOptions(path)
		}
	}

	if dir, ok := os.LookupEnv("HOME"); ok {
		path = filepath.Join(dir, ".config/kmscan/kmscanrc.toml")

		if fileExists(path) {
			return parseOptions(path)
		}

	}

	return parseDefault(kmscanrcExample)
}

func (o *OptionsRC) GetProfile(profileName string) (*Profile, error) {
	if len(o.Profiles) == 0 {
		return nil, fmt.Errorf("at least one profile needs to be defined")
	}

	if profileName == "" {
		return &o.Profiles[0], nil
	}

	for _, p := range o.Profiles {
		if p.Name == profileName {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("profile %s not found", profileName)
}

func parseDefault(kmscanrcExample string) (*OptionsRC, error) {
	opts := &OptionsRC{}

	err := toml.Unmarshal([]byte(kmscanrcExample), &opts)
	if err != nil {
		return opts, fmt.Errorf("unable to parse default config file: %s", err)
	}

	return opts, nil
}

func parseOptions(path string) (*OptionsRC, error) {
	opts := &OptionsRC{}

	bytes, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return opts, fmt.Errorf("invalid path to the config file")
	} else if err != nil {
		return opts, fmt.Errorf("unable to read config file %s: %s", path, err)
	}

	err = toml.Unmarshal(bytes, &opts)
	if err != nil {
		return opts, fmt.Errorf("unable to parse default config file: %s", err)
	}

	return opts, nil
}

func fileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
