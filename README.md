# Estudo API Rest em Golang com Gin
<br> 
<img src="./images/gin_mediun.png">

[<img src="./images/icons/go.svg" width="25px" height="25px" alt="go" title="Go"> <img src="./images/icons/docker.svg" width="25px" height="25px" alt="Docker" title="Docker"> <img src="./images/icons/dotenv.svg" width="25px" height="25px" alt="DotEnv" title="DotEnv"> <img src="./images/icons/github.svg" width="25px" height="25px" alt="GitHub" title="GitHub"> <img src="./images/icons/visualstudiocode.svg" width="25px" height="25px" alt="vscode" title="vscode"> <img src="./images/icons/postgresql.svg" width="25px" height="25px" alt="Postgres" title="Postgres"> <img src="./images/icons/swagger.svg" width="25px" height="25px" alt="Swagger" title="Swagger"> <img src="./images/icons/gatling.svg" width="25px" height="25px" alt="Gatling" title="Gatling">](#estudo-de-autenticação-testes-e-segurança-em-nodejs) <!-- icons by https://simpleicons.org/?q=types -->
<!-- <img src="./images/icons/gatling.svg" width="25px" height="25px" alt="Gatling" title="Gatling"> <img src="./images/icons/githubactions.svg" width="25px" height="25px" alt="Githubactions" title="Githubactions"> <img src="./images/icons/redis.svg" width="25px" height="25px" alt="Redis" title="Redis"> -->



![Badge Status](https://img.shields.io/badge/STATUS-EM_DESENVOLVIMENTO-green)

---

<a id="indice"></a>
## :arrow_heading_up: Índice
<!--ts-->
- [Go: Go e Gin: criando API rest com simplicidade](#estudo-api-rest-em-golang-com-gin)
  :arrow_heading_up: [Índice](#arrow_heading_up-índice)<br/>
  :green_book: [Sobre](#green_book-sobre)<br/>
  :computer: [Rodando o Projeto](#computer-rodando-o-projeto)<br/>
  :newspaper: [Gerando documentação com swagger](#newspaper-gerando-documentação-com-swagger)<br/>
  :camera: [Imagens do Projeto](#camera-imagens-do-projeto)<br/>
  :bar_chart: [Diagramas](#bar_chart-diagramas)<br/>
  :train: [Teste de Carga](#train-teste-de-carga)<br>
  :hammer: [Ferramentas](#hammer-ferramentas)<br/>
  :clap: [Boas Práticas](#clap-boas-práticas)<br/>
  :1234: [Versões](#1234-versões)<br/>

<!--te-->
---
<a id="sobre"></a>
## :green_book: Sobre
Projeto de estudo baseado na trilha [Go e Gin: criando API rest com simplicidade](https://www.alura.com.br/curso-online-go-gin-api-rest-simplicidade). Esse projeto tem finalidade puramente didática. Após a conclusão do projeto do curso, continuei adicionando padrões de mercado como melhorias para estudar algumas aplicações.

As versões mais recentes da linguagem já têm a instalação simplificada pelo `snap`
```bash
$ sudo snap install go --classic
```

Recomendo a instalação do [GVM](https://github.com/moovweb/gvm) para controle de versões da linguagem

Recomendo a instalação da extensão [Golang do VsCode](https://marketplace.visualstudio.com/items?itemName=golang.go)


Descobrindo o host do banco postgres para configurar o pgadmin, apos subir o docker-compose:

```bash
$ docker-compose exec postgres sh
# hostname -i
```
ou
```bash
$ docker inspect container_id | grep IPAddress
```

[:arrow_heading_up: voltar](#indice)

---

### :computer: Rodando o Projeto

Renomeie crie uma copia do arquivo `sample.env` com o nome `.env` e rode o comando docker-compose (de acordo com sua versao do `docker compose`):
```bash
$ docker compose up
```
Aguarde até que as imagens sejam criadas e acesse:

Rota de **headness**: `http://localhost:8080/readiness`
Rota de **liveness**: `http://localhost:8080/liveness`

Acesse para **API**: `http://localhost:8080/alunos`
Acesse para **documentação Swagger**: `http://localhost:8080/swagger/index.html`



[:arrow_heading_up: voltar](#indice)

---
### :newspaper: Gerando documentação com swagger
Para os desenvolvedores que irão manipular o código ou se inspirar para seus próprios desenvolvimentos, há uma particularidade na documentação Swagger. O comando padrão do [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger) (uma ferramenta que gera documentação Swagger para Go) não consegue ler `structs` que utilizam `gorm.Model`, e isso não está explicitamente mencionado em sua documentação. Pesquisando por uma solução, [encontrei o comando apropriado](https://github.com/swaggo/swag/issues/810) para a geração, que segue abaixo:

```bash
$ swag init --parseDependency --parseInternal
```

[:arrow_heading_up: voltar](#indice)

---


### :camera: Imagens do Projeto

<details>
  <summary>Swagger 1</summary>
    <img src="images/captures/swagger_2.png">
</details>
<br>
<details>
  <summary>Swagger 2</summary>
  <img src="images/captures/swagger_2.png">
</details>
<br>
<details>
  <summary>API</summary>
    <img src="images/captures/api.png">
</details>
<br>

[:arrow_heading_up: voltar](#indice)

---

 ### :bar_chart: Diagramas

```mermaid
graph LR
  subgraph Ações Admin Alunos
    A[[ADMIN User]]
    B["Obtém a lista completa de alunos"]
    C["Cria novo aluno"]
    D["Busca aluno por id"]
    E["Deleta aluno por id"]
    F["Edita aluno por id"]
    G["Busca aluno por CPF"]
  end

  subgraph Backend 
    subgraph API
      Aluno["Aluno"]
    end

    subgraph DATABASE
      Aluno-DB[("Aluno-DB")]
    end    
  end

  A --> B
  A --> C
  A --> D
  A --> E
  A --> F
  A --> G

  B -->|GET| Aluno
  C -->|POST| Aluno
  D -->|GET| Aluno
  E -->|DELETE| Aluno
  F -->|PATCH| Aluno
  G -->|GET| Aluno

  Aluno --> Aluno-DB
```
<br>

[:arrow_heading_up: voltar](#indice)

---

### :train: Teste de Carga
**Gatling**
Com o projeto instalado e em execução após o comando `docker compose up`, acesse a rota que renderiza o resultado do teste mais recente em `http://localhost:8082`. Caso você tenha acabado de iniciar o ambiente, nenhum teste terá ocorrido até o momento.

A imagem responsável por fornecer essa saída também é responsável por processar o teste. Para executar um novo teste, basta abrir um novo terminal e, estando na raiz do projeto, execute o comando:

```bash
docker exec -ti gatling-api-test /entrypoint run-test
``` 
Aguarde alguns segundos para o aquecimento dos testes (que inclui o download de dependências, caso não existam, e a execução dos próprios testes). Assim que os testes forem concluídos, o endpoint `http://localhost:8082` apresentará os resultados.

Toda vez que desejar executar os testes novamente, basta rodar o comando a seguir: `docker exec -ti gatling-api-test /entrypoint run-test`.

Estrutura da pasta de testes do Gatling:
```shell
  $ tree
  api-gin-rest
  └── tests
  |    └── gatling
  |        ├── APISimulation  # Pasta de roteiro de testes
  |        |   └── Alunos.java # Roteiro de testes em Java
  |        ├── bundle # Binarios e arquivos Gatling
  |        |   ├── .keep
  |        |   └── ... # Diretórios e arquivos gatling instalados após primeiro teste
  |        └── results # Resultados dos testes
  |        |   ├── latest # Arquivos do ultimo teste performado
  |        |   |   ├── .keep
  |        |   |   └── ... # Diretórios e arquivos do resultado mais recente após primeiro teste
  |        |   └── history # Histórico com todos os testes já performados e a pagina default
  |        |       ├── default # Dados de teste padrão, exibidos quando nenhum teste ainda foi performado
  |        |       └── ... # Diretórios de testes já performados
  |        ├── Dockerfile
  |        └── entrypoin.sh # As automações do Gatling estão aqui.
  |
  $ tree
  .
```
[:arrow_heading_up: voltar](#indice)

---

<a id="ferramentas"></a>
## :hammer: Ferramentas
As seguintes ferramentas foram usadas na construção do projeto:

- [GVM v1.0.22](https://github.com/moovweb/gvm)
- [Go v1.21.1](https://go.dev/)
- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/index.html)
- [Viper](https://github.com/spf13/viper)
- [Gin-Swagger](https://github.com/swaggo/gin-swagger)
- [Postgres](https://www.postgresql.org/)
- [Docker 24.0.6](https://www.docker.com/)
- [Docker compose v2.21.0](https://www.docker.com/)
- [Gatling](https://gatling.io/)
- [VsCode](https://code.visualstudio.com/)
- [DBeaver](https://dbeaver.io/)


[:arrow_heading_up: voltar](#indice)

---


<a id="boas-praticas"></a>
## :clap: Boas Práticas
Seguindo boas práticas de desenvolvimento:
- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)
- [keep a changelog](https://keepachangelog.com/en/1.0.0/)
- [Swagger](https://swagger.io/)

[:arrow_heading_up: voltar](#indice)

---

<a id="versionamento"></a>
## :1234: Versões
As tags de versões estao sendo criadas manualmente a medida que os estudos avançam com melhorias notáveis no projeto. Cada funcionalidade é desenvolvida em uma branch a parte quando finalizadas é gerada tag e mergeadas em master.


Para obter mais informações, consulte o [Histórico de Versões](./CHANGELOG.md).

[:arrow_heading_up: voltar](#indice)
