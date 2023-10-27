package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jtonynet/api-gin-rest/config"
	"github.com/jtonynet/api-gin-rest/models"

	"github.com/jtonynet/api-gin-rest/internal/database"
	"github.com/jtonynet/api-gin-rest/internal/interfaces"
)

type InfraReturn struct {
	message string
	sumary  string
}

// @BasePath /alunos

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

// @Summary Retorna todos os alunos
// @Description Obtém a lista completa de alunos
// @Tags Alunos
// @Produce json
// @Success 200 {array} models.Aluno
// @Router /alunos [get]
func ExibeTodosAlunos(c *gin.Context) {
	var alunos []models.Aluno
	database.DB.Find(&alunos)

	c.Set("result", alunos)
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

	UUID := c.GetString("UUID")
	if UUID == "" {
		UUID = uuid.New().String()
	}
	aluno.UUID = UUID

	c.Set("model", aluno)
	if cfg.FeatureFlags.PostAlunoAsMessageEnabled {
		return
	}

	err := database.DB.Create(&aluno).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao tentar inserir o aluno",
		})
	}

	c.Set("result", aluno)
	c.JSON(http.StatusOK, gin.H{
		"uuid": aluno.UUID})

}

// @Summary Busca aluno por uuid
// @Description Busca aluno por uuid
// @Tags Aluno
// @Produce json
// @Success 200 {object} models.Aluno
// @Failure 404 {string} Not Found
// @Router /aluno/{uuid} [get]
func BuscaAlunoPorUUID(c *gin.Context) {
	var aluno models.Aluno

	uuid := c.Params.ByName("uuid")
	database.DB.Where(&models.Aluno{UUID: uuid}).First(&aluno)

	if aluno.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"Not Found": "Aluno Nao Encontrado",
		})
		return
	}

	c.Set("result", aluno)
	c.JSON(http.StatusOK, aluno)
}

// @Summary Deleta aluno por uuid
// @Description Deleta aluno por uuid
// @Tags Aluno
// @Produce json
// @Success 200 {string} data
// @Router /aluno/{uuid} [delete]
func DeletaAluno(c *gin.Context) {
	var aluno models.Aluno

	uuid := c.Params.ByName("uuid")
	database.DB.Where(&models.Aluno{UUID: uuid}).First(&aluno)

	if err := database.DB.Delete(&aluno).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Aluno inexistente",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "Aluno deletado com sucesso"})
}

// @Summary Edita aluno por uuid
// @Description Edita aluno por uuid
// @Tags Aluno
// @Produce json
// @Accept json
// @Success 202 {uuid} uuid
// @Failure 400 {string} error
// @Router /aluno/{uuid} [patch]
func EditaAluno(c *gin.Context) {
	var aluno models.Aluno

	uuid := c.Params.ByName("uuid")
	database.DB.Where(&models.Aluno{UUID: uuid}).First(&aluno)

	if err := c.ShouldBindJSON(&aluno); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	database.DB.Model(&aluno).UpdateColumns(aluno)
	c.JSON(http.StatusOK, gin.H{
		"uuid": aluno.UUID})
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

	currentTime := time.Now()
	timeFormatted := currentTime.Format("15:04:05.000000")
	fmt.Println("CONTROLLER BuscaAlunoPorCPF (HH:MM:SS.mmmuuu):", timeFormatted)

	c.Set("result", aluno)
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
