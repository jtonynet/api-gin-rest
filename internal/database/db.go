package database

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/api-gin-rest/config"
	"github.com/jtonynet/api-gin-rest/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init(cfg config.Database) error {

	strConn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Pass,
		cfg.DB,
		cfg.Port)

	DB, err = gorm.Open(postgres.Open(strConn))
	if err != nil {
		return err
	}

	DB.AutoMigrate(&models.Aluno{})

	return nil
}

func CheckReadiness() error {
	if err := DB.Raw("SELECT 1").Error; err != nil {
		return err
	}
	return nil
}

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		limit := c.MustGet("Limit").(int64)
		page := c.MustGet("Page").(int64)

		offset := (page - 1) * limit

		return db.Offset(int(offset)).Limit(int(limit))
	}
}
