package lib

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	OdooHome     string            `toml:"odoo_home"`
	DBPrefix     string            `toml:"db_prefix"`
	OdooPort     int               `toml:"odoo_port"`
	Repositories map[string]string `toml:"repositories"`
}

func getDefaultConfig() Config {
	return Config{
		OdooHome: "$ODOO_HOME",
		DBPrefix: "rd-",
		OdooPort: 8069,
		Repositories: map[string]string{
			".workspace": ".workspace",
			"community":  "community",
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

		cfg.OdooHome = os.ExpandEnv(cfg.OdooHome)
		if cfg.OdooHome == "" {
			cfg.OdooHome = "."
		}

		userConfig = &cfg
	})
	return userConfig
}
