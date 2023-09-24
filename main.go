package main

import (
	"log"

	"github.com/jtonynet/api-gin-rest/config"
	"github.com/jtonynet/api-gin-rest/database"
	"github.com/jtonynet/api-gin-rest/routes"
)

// @title api-gin-rest
// @version 1.0
// @description Estudo API Rest em Golang com Gin
// @contact.name API GIN Support
// @termsOfService https://github.com/jtonynet/api-gin-rest
// @contact.url https://github.com/jtonynet/api-gin-rest
// @contact.email learningingenuity@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /alunos
// @Schemes http
// @query.collection.format multi
func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	database.ConectaComBancoDeDados(cfg.Database)
	routes.HandleRequests(cfg.API)
}
