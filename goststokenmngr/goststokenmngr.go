package goststokenmngr

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/reiiissamuel/goststokenmngr/internal"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

func main() {
	// Inicie o scheduler em uma goroutine
	go internal.StartScheduler()

	// Mantenha a aplicação em execução (pode ser um loop infinito, por exemplo)
	select {}
}
