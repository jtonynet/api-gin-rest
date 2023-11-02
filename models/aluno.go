package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
}

type Aluno struct {
	BaseModel `swaggerignore:"true"`
	UUID      uuid.UUID `json:"uuid" example:"db047cc5-193a-4989-93f7-08b81c83eea0" type:uuid;unique`
	Nome      string    `json:"nome" binding:"required" example:"Jonh Doe"`
	CPF       string    `json:"cpf" binding:"required" example:"00000000000"`   //para fins de teste de carga, não valido e nem considero esse campo como unique
	RG        string    `json:"rg" binding:"required" example:"12345678901234"` //para fins de teste de carga, não valido e nem considero esse campo como unique
}
