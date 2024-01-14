package pkg

import "github.com/reiiissamuel/go-sts-token-mngr/internal"

// Função pública para obter o token
func GetCurrentToken() string {
	return internal.GetToken()
}
