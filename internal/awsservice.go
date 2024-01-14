package internal

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

var (
	tokenMutex      sync.Mutex
	awsConfig       *aws.Config
	token           string
	assumeRoleInput sts.AssumeRoleInput
	roleArn         = os.Getenv("AWS_ROLE_ARN")
	roleSessionName = os.Getenv("AWS_ROLE_SESSION_NAME")
	durationSeconds = os.Getenv("AWS_TOKEN_VALID_SECONDS")
	region          = os.Getenv("AWS_REGION")
	stsInstance     *sts.STS
)

const (
	MSG_STARTING_CONFIG  = "Iniciando configurações...\n"
	MSG_CONFIG_COMPLETED = "Configurações concluídas.\n"
	MSG_ERROR            = "Error: %v\n"
)

func GetToken() string {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	return token
}

func StartScheduler() {
	startAWSStAPICallerConfig()

	interval, parseErr := strconv.ParseInt(durationSeconds, 10, 64)
	if parseErr != nil {
		fmt.Printf(MSG_ERROR, parseErr)
	}

	for {
		// Execute a função de busca do token
		if err := updateToken(); err != nil {
			fmt.Printf(MSG_ERROR, err)
		}

		// Aguarde o próximo intervalo
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func startAWSStAPICallerConfig() {
	fmt.Printf(MSG_STARTING_CONFIG)
	awsConfig = &aws.Config{
		Region: aws.String(region), // Defina a região desejada
	}

	getNewSTSInstance()
	getAssumeRoleInput()

	fmt.Printf(MSG_CONFIG_COMPLETED)
}

func updateToken() error {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	// Chame AssumeRole para obter credenciais temporárias
	result, errAssumingRole := stsInstance.AssumeRole(&assumeRoleInput)
	if errAssumingRole != nil {
		return errAssumingRole
	}

	// Atualize o token em memória
	token = *result.Credentials.SessionToken

	//fmt.Println("Token atualizado com sucesso:", token)
	return nil
}

func getNewSTSInstance() {
	// Configuração do AWS STS
	session, err := session.NewSession(awsConfig)
	if err != nil {
		fmt.Printf(MSG_ERROR, err)
	}
	stsInstance = sts.New(session)
}

func getAssumeRoleInput() {
	interval, err := strconv.ParseInt(durationSeconds, 10, 64)

	if err != nil {
		fmt.Printf(MSG_ERROR, err)
	}

	assumeRoleInput = sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(roleSessionName),
		DurationSeconds: aws.Int64(interval),
	}
}
