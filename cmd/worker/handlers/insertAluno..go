package handlers

import (
	"encoding/json"
	"log/slog"

	"github.com/jtonynet/api-gin-rest/internal/database"
	"github.com/jtonynet/api-gin-rest/models"
)

func InsertAluno(msg string) (string, error) {
	var aluno models.Aluno
	err := json.Unmarshal([]byte(msg), &aluno)
	if err != nil {
		return "", err
	}

	err = database.DB.Create(&aluno).Error
	if err != nil {
		slog.Error("cmd:worker:handler:InsertAluno:database.DB.Create error %v", err)
		return "", err
	}

	alunoJSON, err := json.Marshal(aluno)
	if err != nil {
		return "", err
	}

	return string(alunoJSON), nil
}
