package database

import (
	"fmt"
	"log"

	"github.com/jtonynet/api-gin-rest/config"
	"github.com/jtonynet/api-gin-rest/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func ConectaComBancoDeDados(cfg config.Database) {

	strConn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Pass,
		cfg.DB,
		cfg.Port)

	DB, err = gorm.Open(postgres.Open(strConn))
	if err != nil {
		log.Panic("Erro ao conectar com o banco de dados")
	}

	DB.AutoMigrate(&models.Aluno{})
}

func CheckReadiness() error {
	if err := DB.Raw("SELECT 1").Error; err != nil {
		return err
	}
	return nil
}
