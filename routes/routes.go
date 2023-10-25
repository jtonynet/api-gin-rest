package routes

import (
	"fmt"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/jtonynet/api-gin-rest/config"
	"github.com/jtonynet/api-gin-rest/controllers"
	"github.com/jtonynet/api-gin-rest/docs"
	"github.com/jtonynet/api-gin-rest/internal/interfaces"
	"github.com/jtonynet/api-gin-rest/internal/middlewares"
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

	apiGroup.GET("/aluno/:uuid", middlewares.CachedRequest(), controllers.BuscaAlunoPorUUID)
	apiGroup.GET("/aluno/cpf/:cpf", middlewares.CachedRequest(), controllers.BuscaAlunoPorCPF)

	apiGroup.POST("/aluno", controllers.CriaNovoAluno)
	apiGroup.DELETE("/aluno/:id", controllers.DeletaAluno)
	apiGroup.PATCH("/aluno/:id", controllers.EditaAluno)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := fmt.Sprintf(":%s", cfg.Port)
	r.Run(port)
}
