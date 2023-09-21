package main

import (
	"github.com/jtonynet/api-gin-rest/models"
	"github.com/jtonynet/api-gin-rest/routes"
)

func main() {
	models.Alunos = []models.Aluno{
		{Nome: "Jonh Doe", CPF: "06335766889", RG: "47000000000"},
		{Nome: "Marie Doe", CPF: "06335766880", RG: "48000000000"},
	}
	routes.HandleRequests()
}
