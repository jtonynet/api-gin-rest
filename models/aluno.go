package models

import "gorm.io/gorm"

type BaseModel struct {
	gorm.Model
}

type Aluno struct {
	BaseModel `swaggerignore:"true"`
	UUID      string `json:"uuid" example:"00000000-a0a0-0aa0-0aa0-a0aa0000a000" gorm:"unique"`
	Nome      string `json:"nome" binding:"required" example:"Jonh Doe"`
	CPF       string `json:"cpf" binding:"required" example:"00000000000"`   //para fins de teste de carga, não valido e nem considero esse campo como unique
	RG        string `json:"rg" binding:"required" example:"12345678901234"` //para fins de teste de carga, não valido e nem considero esse campo como unique
}
