package main

import (
	"fmt"

	"github.com/sohWenMing/aws_server/internal/awsconnection"
)

func main() {
	cfg, err := awsconnection.InitAWSConfig()
	if err != nil {
		panic(err)
	}
	secret, err := awsconnection.GetSecretFromSecretsManager(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println("returned secret")
	fmt.Printf("%v\n", secret)
	dbString, err := secret.BuildDBString()
	if err != nil {
		panic(err)
	}
	fmt.Printf("db string:\n%s\n", dbString)
}
