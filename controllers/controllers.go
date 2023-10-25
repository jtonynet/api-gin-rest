package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jtonynet/api-gin-rest/config"
	"github.com/jtonynet/api-gin-rest/internal/database"

	"github.com/jtonynet/api-gin-rest/internal/interfaces"
	"github.com/jtonynet/api-gin-rest/models"
)

func Liveness(c *gin.Context) {
	cfg := c.MustGet("cfg").(config.API)

	sumaryData := fmt.Sprintf("%s:%s in TagVersion: %s responds OK",
		cfg.Name,
		cfg.Port,
		cfg.TagVersion)
	c.JSON(http.StatusOK, gin.H{
		"message": "OK", "sumary": sumaryData})
}

func Readiness(c *gin.Context) {
	cfg := c.MustGet("cfg").(config.API)

	var err error

	if err = database.CheckReadiness(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Database Service unavailable",
		})
		return
	}

	if cfg.FeatureFlags.CacheEnabled {
		cacheClient := c.MustGet("cacheClient").(interfaces.CacheClient)
		if !cacheClient.IsConnected() {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"message": "CacheClient Service unavailable",
			})
			return
		}
	}

	if cfg.FeatureFlags.PostAlunoAsMessageEnabled {
		messageBroker := c.MustGet("messageBroker").(interfaces.Broker)
		if !messageBroker.IsConnected() {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"message": "MessageBroker Service unavailable",
			})
			return
		}
	}

	sumaryData := fmt.Sprintf("%s:%s in TagVersion: %s responds: OK",
		cfg.Name,
		cfg.Port,
		cfg.TagVersion)

	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"sumary":  sumaryData,
	})
}

// @BasePath /alunos

// @Summary Retorna todos os alunos
// @Description Obtém a lista completa de alunos
// @Tags Alunos
// @Produce json
// @Success 200 {array} models.Aluno
// @Router /alunos [get]
func ExibeTodosAlunos(c *gin.Context) {
	var alunos []models.Aluno
	database.DB.Find(&alunos)
	c.JSON(200, alunos)
}

// @Summary Cria novo aluno
// @Description Cria novo aluno
// @Tags Aluno
// @Produce json
// @Accept json
// @Success 200 {object} models.Aluno
// @Success 202 {uuid} uuid
// @Failure 400 {string} Error
// @Router /aluno [post]
func CriaNovoAluno(c *gin.Context) {
	cfg := c.MustGet("cfg").(config.API)

	var aluno models.Aluno
	if err := c.ShouldBindJSON(&aluno); err != nil {
		slog.Error("controllers:CriaNovoAluno:ShouldBindJSON error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	aluno.UUID = uuid.New().String()
	var cacheClient interfaces.CacheClient

	if cfg.FeatureFlags.CacheEnabled {
		cacheClient = c.MustGet("cacheClient").(interfaces.CacheClient)
		var err error

		var msgJson map[string]interface{}
		msg := fmt.Sprintf(`{"HTTPStatusCode":202, "uuid":"%s", "Message":" in processing."}`, aluno.UUID)
		err = json.Unmarshal([]byte(msg), &msgJson)
		if err != nil {
			slog.Error("controllers:CriaNovoAluno:json.Unmarshal error: %v", err)
		}

		expiration := time.Duration(0 * time.Millisecond) // Request Post no expires INFO cache. Only DATA cache expires
		err = cacheClient.Set(aluno.UUID, msg, expiration)
		if err != nil {
			slog.Error("controllers:CriaNovoAluno:cacheClient.Set error: %v", err)
		}
	}

	if cfg.FeatureFlags.PostAlunoAsMessageEnabled {
		messageBroker := c.MustGet("messageBroker").(interfaces.Broker)

		alunoJSON, err := json.Marshal(aluno)
		if err != nil {
			slog.Error("controllers:CriaNovoAluno:json.Marshal error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao serializar Aluno em JSON",
			})

			return
		}

		err = messageBroker.Publish(string(alunoJSON))
		if err != nil {
			slog.Error("controllers:CriaNovoAluno:messageBroker.Publish error: %v", err)
		}

		c.JSON(http.StatusAccepted, gin.H{
			"uuid": aluno.UUID})

		return
	} else {
		err := database.DB.Create(&aluno).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro ao tentar inserir o aluno",
			})
		}

		if cfg.FeatureFlags.CacheEnabled {
			alunoJSON, err := json.Marshal(aluno)
			if err != nil {
				slog.Error("controllers:CriaNovoAluno:json.Marshal error %v", err)
			}

			err = cacheClient.Set(aluno.UUID, string(alunoJSON), cacheClient.GetDefaultExpiration())
			if err != nil {
				slog.Error("controllers:CriaNovoAluno:cacheClient.Set error set%v", err)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"uuid": aluno.UUID})

		return
	}
}

// @Summary Busca aluno por id
// @Description Busca aluno por id
// @Tags Aluno
// @Produce json
// @Success 200 {object} models.Aluno
// @Failure 404 {string} Not Found
// @Router /aluno/{id} [get]
func BuscaAlunoPorId(c *gin.Context) {
	var aluno models.Aluno
	id := c.Params.ByName("id")
	database.DB.First(&aluno, id)

	if aluno.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Aluno Nao Encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, aluno)
}

// @Summary Busca aluno por id
// @Description Busca aluno por id
// @Tags Aluno
// @Produce json
// @Success 200 {object} models.Aluno
// @Failure 404 {string} Not Found
// @Router /aluno/uuid/{uuid} [get]
func BuscaAlunoPorUUID(c *gin.Context) {
	// cfg := c.MustGet("cfg").(config.API)

	var aluno models.Aluno
	uuid := c.Params.ByName("uuid")

	database.DB.Where(&models.Aluno{UUID: uuid}).First(&aluno)

	if aluno.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Aluno Nao Encontrado",
		})
		return
	}

	// var cacheClient interfaces.CacheClient
	// if cfg.FeatureFlags.CacheEnabled {
	// 	cacheClient = c.MustGet("cacheClient").(interfaces.CacheClient)

	// 	alunoJSON, err := json.Marshal(aluno)
	// 	if err != nil {
	// 		slog.Error("controllers:BuscaAlunoPorUUID:json.Marshal error: %v", err)
	// 	}

	// 	err = cacheClient.Set(aluno.UUID, string(alunoJSON), cacheClient.GetDefaultExpiration())
	// 	if err != nil {
	// 		slog.Error("controllers:BuscaAlunoPorUUID:cacheClient.Set error: %v", err)
	// 	}
	// }

	currentTime := time.Now()
	timeFormatted := currentTime.Format("15:04:05.000000")
	fmt.Println("2 - RETORNANDO DA ROTA CONTROLLER (HH:MM:SS.mmmuuu):", timeFormatted)

	alunoJSON, err := json.MarshalIndent(aluno, "", "  ")
	if err != nil {
		// Lida com erros se a conversão falhar
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro na conversão para JSON",
		})
		return
	}

	var returnData map[string]interface{}
	if err := json.Unmarshal([]byte(alunoJSON), &returnData); err != nil {
		slog.Error("middlewares:CachedRequest:json.Unmarshal error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao recuperar dados",
		})
		c.Abort()
	}

	c.Set("queryResult", string(alunoJSON))

	c.JSON(http.StatusOK, aluno)
}

// @Summary Deleta aluno por id
// @Description Deleta aluno por id
// @Tags Aluno
// @Produce json
// @Success 200 {string} data
// @Router /aluno/{id} [delete]
func DeletaAluno(c *gin.Context) {
	var aluno models.Aluno
	id := c.Params.ByName("id")
	database.DB.Delete(&aluno, id)
	c.JSON(http.StatusOK, gin.H{
		"data": "Aluno deletado com sucesso"})
}

// @Summary Edita aluno por id
// @Description Edita aluno por id
// @Tags Aluno
// @Produce json
// @Accept json
// @Success 200 {object} models.Aluno
// @Failure 400 {string} error
// @Router /aluno/{id} [patch]
func EditaAluno(c *gin.Context) {
	var aluno models.Aluno
	id := c.Params.ByName("id")
	database.DB.First(&aluno, id)

	if err := c.ShouldBindJSON(&aluno); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	database.DB.Model(&aluno).UpdateColumns(aluno)
	c.JSON(http.StatusOK, aluno)
}

// @Summary Busca aluno por CPF
// @Description Busca aluno por CPF
// @Tags Aluno
// @Produce json
// @Success 200 {object} models.Aluno
// @Failure 404 {string} Not Found
// @Router /aluno/cpf/{cpf} [get]
func BuscaAlunoPorCPF(c *gin.Context) {
	var aluno models.Aluno
	cpf := c.Param("cpf")

	database.DB.Where(&models.Aluno{CPF: cpf}).First(&aluno)

	if aluno.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Aluno não Encontrado"})
		return
	}

	c.JSON(http.StatusOK, aluno)
}

// Importe as bibliotecas e pacotes necessários aqui

// @Summary Busca o número de alunos no banco de dados
// @Description Busca o número de alunos no banco de dados
// @Tags Aluno
// @Produce json
// @Success 200 {object} int
// @Router /alunos/count [get]
func ContaAlunos(c *gin.Context) {
	var totalAlunos int64
	database.DB.Model(&models.Aluno{}).Count(&totalAlunos)

	c.JSON(http.StatusOK, totalAlunos)
}
