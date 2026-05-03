package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	Env    string
	VarDir string
	Host   string
	Port   int
}

func Load() *Env {
	_ = godotenv.Load()

	env := &Env{
		Env:    "dev",
		VarDir: "var",
		Host: "127.0.0.1",
		Port: 3096,
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
		portInt, err := strconv.Atoi(port)
		if err != nil {
			portInt = 0
		}

		env.Port = portInt
	}

	return env
}

func (e *Env) Validate() error {
	if e.Env != "dev" && e.Env != "prod" {
		return errors.New("invalid env provided (it must be 'dev' or 'prod')")
	}

	if e.Port < 1 || e.Port > 65535 {
		return errors.New("invalid port provided")
	}

	return nil
}

func (e *Env) GetLogsDir() string {
	return fmt.Sprintf("%s/logs", e.VarDir)
}
