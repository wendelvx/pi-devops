package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConectDb_PanicOnConnectionFailure(t *testing.T) {
	// Setup: Definimos variáveis de ambiente inválidas
	os.Setenv("DB_HOST", "host_inexistente")
	os.Setenv("DB_USER", "user_teste")
	os.Setenv("DB_PASS", "pass_teste")
	os.Setenv("DB_NAME", "db_teste")

	// Cleanup: Limpa as variáveis de ambiente ao fim do teste
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("DB_NAME")
	}()

	// Act & Assert: O Testify testa o panic diretamente
	assert.Panics(t, func() {
		ConectDb()
	}, "A função ConectDb deveria ter disparado um panic devido ao banco inexistente")
}

func TestConectDb_FallbackLocalhost(t *testing.T) {
	// Setup: Garante que o host está vazio
	os.Setenv("DB_HOST", "")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASS", "pass")
	os.Setenv("DB_NAME", "db")

	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("DB_NAME")
	}()

	// Act & Assert
	assert.Panics(t, func() {
		ConectDb()
	}, "A função deveria tentar conectar no localhost e falhar (panic)")
}
