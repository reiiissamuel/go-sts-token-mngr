package pkg

import "github.com/reiiissamuel/goststokenmngr/internal"

// Função pública para obter o token
func GetCurrentToken() string {
	return internal.GetToken()
}
