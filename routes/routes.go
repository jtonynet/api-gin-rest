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
	messageBroker interfaces.MessageBroker,
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

	apiGroup.GET("/alunos/count", controllers.ContaAlunos)

	apiGroup.GET("/alunos",
		middlewares.PaginateRequest(),
		middlewares.CachedGetRequest(),
		controllers.ExibeTodosAlunos,
	)

	apiGroup.GET("/aluno/:uuid",
		middlewares.CachedGetRequest(),
		controllers.BuscaAlunoPorUUID,
	)

	apiGroup.GET("/aluno/cpf/:cpf",
		middlewares.CachedGetRequest(),
		controllers.BuscaAlunoPorCPF,
	)

	apiGroup.POST("/aluno",
		middlewares.CachedPostRequest(),
		middlewares.MessageBrokerPublishPostRequest(),
		controllers.CriaNovoAluno,
	)

	apiGroup.DELETE("/aluno/:uuid",
		middlewares.CachedDeleteRequest(),
		controllers.DeletaAluno,
	)
	apiGroup.PATCH("/aluno/:uuid",
		middlewares.CachedDeleteRequest(),
		controllers.EditaAluno,
	)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := fmt.Sprintf(":%s", cfg.Port)
	r.Run(port)
}
