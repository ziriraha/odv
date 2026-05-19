package lib

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

type OdooConfig struct {
	Path    string `toml:"path"`
	Command string `toml:"command"`
}

type DBConfig struct {
	DSN    string `toml:"dsn"`
	Prefix string `toml:"prefix"`
}

type Config struct {
	Odoo         OdooConfig        `toml:"odoo"`
	Database     DBConfig          `toml:"db"`
	Repositories map[string]string `toml:"repositories"`
}

const defaultDSN = "postgres://$PGUSER:$PGPASSWORD@$PGHOST:$PGPORT/postgres"

func getDefaultConfig() Config {
	return Config{
		Odoo: OdooConfig{
			Path:    ".",
			Command: "python3 ./odoo/odoo-bin",
		},
		Database: DBConfig{
			DSN:    defaultDSN,
			Prefix: "odoo",
		},
		Repositories: map[string]string{
			".workspace": ".workspace",
			"odoo":       "odoo",
			"enterprise": "enterprise",
			"upgrade":    "upgrade",
		},
	}
}

var (
	userConfig     *Config
	userConfigOnce sync.Once

	userHome     string
	userHomeOnce sync.Once
)

func GetUserHome() string {
	userHomeOnce.Do(func() {
		var err error
		userHome, err = os.UserHomeDir()
		if err != nil {
			panic("Failed to get user home directory: " + err.Error())
		}
	})
	return userHome
}

func GetConfig() *Config {
	userConfigOnce.Do(func() {
		cfg := getDefaultConfig()

		configPath := filepath.Join(GetUserHome(), ".odvrc")
		data, err := os.ReadFile(configPath)
		if err == nil {
			toml.Unmarshal(data, &cfg)
		}

		cfg.Odoo.Path = os.ExpandEnv(cfg.Odoo.Path)
		cfg.Database.DSN = os.ExpandEnv(cfg.Database.DSN)

		userConfig = &cfg
	})
	return userConfig
}
