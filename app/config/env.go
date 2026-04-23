package config

import (
	"fmt"
	"os"
)

type Env struct {
	Env      string `mapstructure:"ENV" default:"dev"`
	VarDir   string `mapstructure:"VAR_DIR" default:"var"`
	Host     string `mapstructure:"HOST" default:"127.0.0.1"`
	Port     int    `mapstructure:"PORT" default:"8080"`
	LogLevel string `mapstructure:"LOG_LEVEL" default:"info"`
}

func NewEnv() (*Env, error) {
	return Load()
}

func Load() (*Env, error) {
	env := &Env{
		Env:      "dev",
		VarDir:   "var",
		Host:     "127.0.0.1",
		Port:     8080,
		LogLevel: "info",
	}

	if environment := os.Getenv("ENV"); environment != "" {
		env.Env = environment
	}

	if varDir := os.Getenv("VAR_DIR"); varDir != "" {
		env.VarDir = varDir
	}

	if host := os.Getenv("HOST"); host != "" {
		env.Host = host
	}

	if port := os.Getenv("PORT"); port != "" {
		fmt.Sscanf(port, "%d", &env.Port)
	}

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		env.LogLevel = level
	}

	return env, nil
}

func (e *Env) Validate() error {
	return nil
}

func (e *Env) GetLogsDir() string {
	return e.VarDir + "/logs"
}

func (e *Env) GetUploadsDir() string {
	return e.VarDir + "/uploads"
}