package envloader

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/joho/godotenv"
	"github.com/sohWenMing/aws_server/internal/utils"
)

type Environment string

var environments = []Environment{
	"DEV",
	"DEV_TO_PROD",
	"PROD",
}

type prodDbCredentials struct {
	user       string
	password   string
	dbName     string
	sslMode    string
	projectDir string
}

func (p prodDbCredentials) buildDBstring() string {
	return fmt.Sprintf("postgresql://%s:%s@localhost:5432/%s?sslmode=%s&sslrootcert=%s/global-bundle.pem",
		p.user, p.password, p.dbName, p.sslMode, p.projectDir,
	)
}

func getProdDbCredentials() prodDbCredentials {
	returned := prodDbCredentials{}
	returned.user = os.Getenv("PROD_POSTGRES_USER")
	returned.password = os.Getenv("PROD_POSTGRES_PASSWORD")
	returned.dbName = os.Getenv("PROD_DBNAME")
	returned.sslMode = os.Getenv("PROD_SSLMODE")
	projectDir, err := utils.GetProjectRoot()
	if err != nil {
		log.Fatal(err)
	}
	if returned.user == "" || returned.password == "" || returned.dbName == "" || returned.sslMode == "" {
		errMsg := fmt.Sprintf("mandatory value in prodDbCredentials not found. prodDbCredentials:\n%v\n", returned)
		log.Fatal(errMsg)
	}
	returned.projectDir = projectDir
	return returned
}

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
	environment, err := getEnvironmentFromEnvString(env)
	if err != nil {
		return err
	}
	if environment == "DEV" {
		dbString := os.Getenv("DBSTRING")
		if dbString == "" {
			return errors.New("dbstring could not be found for setting in dev environment")
		}
		e.dbString = dbString
		return nil
	}
	if env == "DEV_TO_PROD" {
		e.dbString = getProdDbCredentials().buildDBstring()
		return nil
	}
	if env == "PROD" {
		e.dbString = "TODO"
		return nil
	}
	return errors.New("unrecognized environment in SetDBString")
}

func (e *EnvSettings) SetEnv(env string) error {
	environment, err := getEnvironmentFromEnvString(env)
	if err != nil {
		return err
	}
	e.env = string(environment)
	return nil
}

func getEnvironmentFromEnvString(env string) (environment Environment, err error) {
	environment = Environment(env)
	if !slices.Contains(environments, environment) {
		return "", fmt.Errorf("%q is not a recognized environment\n", environment)
	}
	return environment, nil
}

func LoadEnvironment(envPath string) (envSettings *EnvSettings, err error) {
	envSettings = &EnvSettings{}
	env := os.Getenv("ENV")
	environment, err := getEnvironmentFromEnvString(env)
	if err != nil {
		return nil, err
	}
	if string(environment) == "DEV" || string(environment) == "DEV_TO_PROD" {
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
