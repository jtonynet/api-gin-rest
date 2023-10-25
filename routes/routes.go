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

	"github.com/jtonynet/api-gin-rest/internal/interfaces"
)

func HandleRequests(
	cfg config.API,
	messageBroker interfaces.Broker,
	cacheClient interfaces.CacheClient,
) {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/"

	if cfg.FeatureFlags.PprofCPUEnabled {
		pprof.Register(r, "/debug/pprof")
	}

	apiGroup := r.Group("/")

	apiGroup.Use(middlewares.ConfigInject(cfg))
	apiGroup.Use(middlewares.CacheClientInject(cacheClient))
	apiGroup.Use(middlewares.MessageBrokerInject(messageBroker))

	apiGroup.GET("/liveness", controllers.Liveness)
	apiGroup.GET("/readiness", controllers.Readiness)

	apiGroup.GET("/alunos", controllers.ExibeTodosAlunos)
	apiGroup.GET("/alunos/count", controllers.ContaAlunos)

	apiGroup.GET("/aluno/uuid/:uuid", middlewares.CachedRequest(cacheClient), controllers.BuscaAlunoPorUUID)

	apiGroup.POST("/aluno", controllers.CriaNovoAluno)
	apiGroup.GET("/aluno/:id", controllers.BuscaAlunoPorId)
	apiGroup.DELETE("/aluno/:id", controllers.DeletaAluno)
	apiGroup.PATCH("/aluno/:id", controllers.EditaAluno)
	apiGroup.GET("/aluno/cpf/:cpf", controllers.BuscaAlunoPorCPF)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := fmt.Sprintf(":%s", cfg.Port)
	r.Run(port)
}
