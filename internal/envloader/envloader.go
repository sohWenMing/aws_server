package envloader

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type EnvSettings struct {
	dbString string
	env      string
}

func (e *EnvSettings) GetDBString() string {
	return e.dbString
}
func (e *EnvSettings) GetEnv() string {
	return e.env
}
func (e *EnvSettings) SetDBString(env string) error {
	if env == "DEV" {
		dbString := os.Getenv("DBSTRING")
		if dbString == "" {
			return errors.New("dbstring could not be found for setting in dev environment")
		}
		e.dbString = dbString
		return nil
	}
	if env == "PROD" {
		e.dbString = "TODO"
		return nil
	}
	return errors.New("unrecognized environment in SetDBString")
}
func (e *EnvSettings) SetEnv(env string) error {
	if env != "DEV" && env != "PROD" {
		return fmt.Errorf("%q is not a recognized environment\n", env)
	}
	e.env = env
	return nil
}

func LoadEnvironment(envPath string) (envSettings *EnvSettings, err error) {
	envSettings = &EnvSettings{}
	env := os.Getenv("ENV")
	if env != "DEV" && env != "PROD" {
		return nil, fmt.Errorf("%q is not a recognized environment\n", env)

	}
	if env == "DEV" {
		err := loadEnvFile(envPath)
		if err != nil {
			return nil, err
		}
	}
	err = envSettings.SetDBString(env)
	if err != nil {
		return nil, err
	}
	err = envSettings.SetEnv(env)
	if err != nil {
		return nil, err
	}
	return envSettings, nil
}

func loadEnvFile(envPath string) error {
	err := godotenv.Load(envPath)
	if err != nil {
		return err
	}
	return nil
}
