package database

import (
	"fmt"
	"log"

	"github.com/jtonynet/api-gin-rest/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConectaComBancoDeDados() {
	//stringDeConexao := "host=localhost user=api_user password=api_pass dbname=api_gin_rest_db port=5432 sslmode=disable"

	fmt.Print("Entrei")
	stringDeConexao := "host=postgres-gin-rest user=api_user password=api_pass dbname=api_gin_rest_db port=5432 sslmode=disable"

	DB, err = gorm.Open(postgres.Open(stringDeConexao))
	if err != nil {
		log.Panic("Erro ao conectar com o banco de dados")
	}
	DB.AutoMigrate(&models.Aluno{})
}
