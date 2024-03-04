package config

import (
	"time"
	"upgrader/pkg"
)

// Config app config
type Config struct {
	URL       string        `yaml:"url"`
	Version   string        `yaml:"version"`
	AuthCode  string        `yaml:"authCode"`
	RunnerDir string        `yaml:"runnerDir"`
	BackupDir string        `yaml:"backupDir"`
	TempDir   string        `yaml:"tempDir"`
	FileName  string        `yaml:"fileName"`
	Timer     time.Duration `yaml:"timer"`
}

func CreateEmpty() *Config {
	return &Config{
		URL:       "http://limit.api.yyxcloud.com/",
		Version:   "0",
		AuthCode:  "authCode",
		RunnerDir: "runner",
		BackupDir: "backup",
		TempDir:   "temp",
		FileName:  "ledshow.tar",
		Timer:     5 * time.Minute,
	}
}
func Load(path string) (*Config, error) {
	// Initialize options from config file and CLI context.
	var config Config
	err := pkg.LoadYml(path, &config)
	if err != nil {
		pkg.Log.Error(err)
	}
	return &config, nil
}
func (c *Config) Save(path string) error {
	return pkg.SaveYml("config.yml", c)
}
