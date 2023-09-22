package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/api-gin-rest/controllers"
)

func HandleRequests() {
	r := gin.Default()
	r.GET("/alunos", controllers.ExibeTodosAlunos)
	r.GET("/:nome", controllers.Saudacao)
	r.POST("/alunos", controllers.CriaNovoAluno)
	r.GET("/alunos/:id", controllers.BuscaAlunoPorId)
	r.DELETE("/alunos/:id", controllers.DeletaAluno)
	r.PATCH("/alunos/:id", controllers.EditaAluno)
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)
	r.Run()
}
