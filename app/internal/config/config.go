package config

import (
	logger "emails_extractor/app/pkg/logger/zap"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	Databases struct {
		Input        string `yaml:"input" env-default:"./databases/"`
		Output       string `yaml:"output" env-default:"./results/"`
		MergedOutput string `yaml:"merged_output" env-default:"./all/"`
	} `yaml:"databases"`
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		logger.GetLogger().Info("initialing config...")
		instance = &Config{}

		if err := cleanenv.ReadConfig("./config/config.yml", instance); err != nil {
			description, err := cleanenv.GetDescription(instance, nil)

			logger.GetLogger().Info("failed to read config")
			logger.GetLogger().Warn(description)
			logger.GetLogger().Error(err.Error())
		}
	})

	return instance
}
