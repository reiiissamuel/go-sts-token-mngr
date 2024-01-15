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
	token           string
	assumeRoleInput sts.AssumeRoleInput
	stsInstance     *sts.STS
)

const (
	MSG_STARTING_CONFIG  = "Iniciando configurações..."
	MSG_CONFIG_COMPLETED = "Configurações concluídas."
	MSG_STARTING_JOB     = "Iniciando Job..."
	MSG_ERROR            = "Error: %v\n"
)

func GetToken() string {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	return token
}

func StartScheduler() {
	startAWSStAPICallerConfig()
	interval, parseErr := strconv.ParseInt(os.Getenv("AWS_TOKEN_VALID_SECONDS"), 10, 64)
	if parseErr != nil {
		fmt.Printf(MSG_ERROR, parseErr)
	}
	fmt.Println(MSG_STARTING_JOB)
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
	fmt.Println(MSG_STARTING_CONFIG)

	getNewSTSInstance()
	getAssumeRoleInput()

	fmt.Println(MSG_CONFIG_COMPLETED)
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
	session, err := session.NewSession()
	if err != nil {
		fmt.Printf(MSG_ERROR, err)
	}
	stsInstance = sts.New(session)
}

func getAssumeRoleInput() {
	interval, err := strconv.ParseInt(os.Getenv("AWS_TOKEN_VALID_SECONDS"), 10, 64)

	if err != nil {
		fmt.Printf(MSG_ERROR, err)
	}

	assumeRoleInput = sts.AssumeRoleInput{
		RoleArn:         aws.String(os.Getenv("AWS_ROLE_ARN")),
		RoleSessionName: aws.String(os.Getenv("AWS_ROLE_SESSION_NAME")),
		DurationSeconds: aws.Int64(interval),
	}
}
