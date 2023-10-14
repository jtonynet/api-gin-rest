package routes

import (
	"fmt"

	pprof "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/api-gin-rest/config"
	"github.com/jtonynet/api-gin-rest/controllers"
	docs "github.com/jtonynet/api-gin-rest/docs"
	"github.com/jtonynet/api-gin-rest/internal/middlewares"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func HandleRequests(cfg config.API) {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/"

	if cfg.FeatureFlags.PprofCPUEnabled {
		pprof.Register(r, "/debug/pprof")
	}

	apiGroup := r.Group("/")

	apiGroup.Use(middlewares.ConfigInjectHandler(cfg))
	apiGroup.Use(middlewares.ConfigManagerHandler())

	apiGroup.GET("/liveness", controllers.Liveness)
	apiGroup.GET("/readiness", controllers.Readiness)

	apiGroup.GET("/alunos", controllers.ExibeTodosAlunos)
	apiGroup.GET("/alunos/count", controllers.ContaAlunos)

	apiGroup.GET("/aluno/uuid/:uuid", controllers.BuscaAlunoPorUUID)

	apiGroup.POST("/aluno", controllers.CriaNovoAluno)
	apiGroup.GET("/aluno/:id", controllers.BuscaAlunoPorId)
	apiGroup.DELETE("/aluno/:id", controllers.DeletaAluno)
	apiGroup.PATCH("/aluno/:id", controllers.EditaAluno)
	apiGroup.GET("/aluno/cpf/:cpf", controllers.BuscaAlunoPorCPF)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := fmt.Sprintf(":%s", cfg.Port)
	r.Run(port)
}
