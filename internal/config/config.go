package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ShortLinkLength int      `yaml:"shortLinkLength"`
	Storage         string   `yaml:"storage"`
	Server          Server   `yaml:"server"`
	Database        Database `yaml:"database"`
}

type Server struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type Database struct {
	MigrationsDir string `yaml:"migrationsDir"`
	Host          string `yaml:"host"`
	Port          string `yaml:"port"`
	Username      string `yaml:"user"`
	Password      string `yaml:"pass"`
	DBName        string `yaml:"dbName"`
}

func NewFromFile(log *logrus.Logger, filepath string) (Config, error) {
	var config Config

	f, err := os.Open(filepath)
	if err != nil {
		return config, err
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Warnf("Close config file: %v", err)
		}
	}()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)

	return config, err
}

func ParseConfig(log *logrus.Logger) (Config, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	defaultPath := filepath.Join(dir, "config", "main.yml")
	configFilePath := flag.String("config", defaultPath, "config file path")

	flag.Parse()

	configFileFullPath, err := filepath.Abs(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := NewFromFile(log, configFileFullPath)
	if err != nil {
		return cfg, fmt.Errorf("unable to load config: %w", err)
	}

	return cfg, nil
}
