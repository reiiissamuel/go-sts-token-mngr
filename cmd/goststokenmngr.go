package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/reiiissamuel/go-sts-token-mngr/internal"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
}

func main() {
	// Inicie o scheduler em uma goroutine
	go internal.StartScheduler()

	// Mantenha a aplicação em execução (pode ser um loop infinito, por exemplo)
	select {}
}
