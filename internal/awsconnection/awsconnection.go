package awsconnection

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/sohWenMing/aws_server/internal/utils"
)

type DBSecret struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	Engine               string `json:"engine"`
	Host                 string `json:"host"`
	Port                 int    `json:"port"`
	DBName               string `json:"dbname"`
	DBInstanceIdentifier string `json:"dbInstanceIdentifier"`
}

func (p DBSecret) BuildDBString() (dbString string, err error) {
	projectRoot, err := utils.GetProjectRoot()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=verify-full&sslrootcert=%s/global-bundle.pem",
		p.Username,
		p.Password,
		p.Host,
		p.Port,
		p.DBName,
		projectRoot,
	), nil
}

const region = "ap-southeast-1"

func InitAWSConfig() (returned aws.Config, err error) {
	nilCfg := aws.Config{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	returned, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	select {
	case <-ctx.Done():
		return nilCfg, ctx.Err()
	default:
		if err != nil {
			return nilCfg, err
		}
		return returned, nil
	}
}

func GetSecretFromSecretsManager(cfg aws.Config) (dbSecret DBSecret, err error) {
	nilSecret := DBSecret{}
	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("nindgabeet/db-credentials"),
		VersionStage: aws.String("AWSCURRENT"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := svc.GetSecretValue(ctx, input)
	select {
	case <-ctx.Done():
		return nilSecret, ctx.Err()
	default:
		if err != nil {
			return nilSecret, err
		}
		secretString := []byte(*result.SecretString)
		var returned DBSecret
		if err := json.Unmarshal(secretString, &returned); err != nil {
			return nilSecret, err
		}
		return returned, nil

	}
}
