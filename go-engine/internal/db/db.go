package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConectDb() *gorm.DB {
	// Lendo as variáveis de ambiente que definimos no docker-compose
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	// Fallback para localhost caso você rode o binário fora do Docker para testes
	if host == "" {
		host = "localhost"
	}

	// Montando a DSN (Data Source Name) dinamicamente
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=America/Sao_Paulo",
		host, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		// Panic aqui é útil para o Docker saber que o container falhou e tentar o restart
		panic(fmt.Sprintf("Erro ao conectar no banco de dados: %v", err))
	}

	return db
}