package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/api-gin-rest/config"
	"github.com/jtonynet/api-gin-rest/controllers"
	docs "github.com/jtonynet/api-gin-rest/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func HandleRequests(cfg config.API) {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/"

	r.GET("/liveness", func(c *gin.Context) {
		controllers.Liveness(c, cfg)
	})

	r.GET("/headness", func(c *gin.Context) {
		controllers.Headness(c, cfg)
	})

	r.GET("/alunos", controllers.ExibeTodosAlunos)
	r.POST("/alunos", controllers.CriaNovoAluno)
	r.GET("/alunos/:id", controllers.BuscaAlunoPorId)
	r.DELETE("/alunos/:id", controllers.DeletaAluno)
	r.PATCH("/alunos/:id", controllers.EditaAluno)
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(cfg.Port)
}
