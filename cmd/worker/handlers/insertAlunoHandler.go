package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/jtonynet/api-gin-rest/internal/database"
	"github.com/jtonynet/api-gin-rest/models"
)

func insertAlunoHandler(msg string) error {
    var aluno models.Aluno
    err := json.Unmarshal([]byte(msg), &aluno)
    if err != nil {
        // fmt.Println("REQUEUE: Erro na análise JSON: ", err)
        return err
    }

    err = database.DB.Create(&aluno).Error
    if err != nil {
        // fmt.Println("REQUEUE: Erro no insert do BD: ", err)
        return err
    }

    return nil
}