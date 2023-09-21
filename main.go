package main

import (
	"github.com/jtonynet/api-gin-rest/database"
	"github.com/jtonynet/api-gin-rest/routes"
)

func main() {
	database.ConectaComBancoDeDados()
	routes.HandleRequests()
}
