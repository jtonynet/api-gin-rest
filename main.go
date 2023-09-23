package main

import (
	"fmt"

	"github.com/jtonynet/api-gin-rest/database"
	"github.com/jtonynet/api-gin-rest/routes"
)

func main() {
	database.ConectaComBancoDeDados()
	routes.HandleRequests()
	fmt.Println("Olar Mundo")
}

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// type Personalidade struct {
// 	Nome     string `json:"nome"`
// 	Historia string `json:"historia"`
// }

// func main() {
// 	//database.ConectaComBancoDeDados()
// 	http.HandleFunc("/", Home)
// 	http.HandleFunc("/api/personalidades", TodasPersonalidades)

// 	fmt.Println("Iniciando o server Rest com go")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func Home(w http.ResponseWriter, _ *http.Request) {
// 	fmt.Fprint(w, "Home Page")
// }

// func TodasPersonalidades(w http.ResponseWriter, _ *http.Request) {
// 	Personalidades := []Personalidade{
// 		{Nome: "Nome1", Historia: "Historia 1"},
// 		{Nome: "Nome2", Historia: "Historia 2"},
// 	}

// 	json.NewEncoder(w).Encode(Personalidades)
// }
