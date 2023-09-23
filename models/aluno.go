package models

import "gorm.io/gorm"

type BaseModel struct {
	gorm.Model
}

type Aluno struct {
	BaseModel `swaggerignore:"true"`
	Nome      string `json:"nome" example:"Jonh Doe"`
	CPF       string `json:"cpf" example:"00000000000"`
	RG        string `json:"rg" example:"12345678901234"`
}
