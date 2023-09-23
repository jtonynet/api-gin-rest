# Estudo API Rest em Golang com Gin
<br> 
<img src="./images/gopher_2.png">

[<img src="./images/icons/go.svg" width="25px" height="25px" alt="go" title="Go"> <img src="./images/icons/docker.svg" width="25px" height="25px" alt="Docker" title="Docker"> <img src="./images/icons/github.svg" width="25px" height="25px" alt="GitHub" title="GitHub"> <img src="./images/icons/visualstudiocode.svg" width="25px" height="25px" alt="vscode" title="vscode"> <img src="./images/icons/postgresql.svg" width="25px" height="25px" alt="Postgres" title="Postgres"> <img src="./images/icons/swagger.svg" width="25px" height="25px" alt="Swagger" title="Swagger">](#estudo-de-autenticação-testes-e-segurança-em-nodejs) <!-- icons by https://simpleicons.org/?q=types -->
<!-- <img src="./images/icons/gatling.svg" width="25px" height="25px" alt="Gatling" title="Gatling"> -->



![Badge Status](https://img.shields.io/badge/STATUS-EM_DESENVOLVIMENTO-green)

---

<a id="indice"></a>
## :arrow_heading_up: Índice
<!--ts-->
- [Go: Go e Gin: criando API rest com simplicidade](#estudo-api-rest-em-golang-com-gin)<br>
  :arrow_heading_up: [Índice](#arrow_heading_up-índice)<br>
  :green_book: [Sobre](#green_book-sobre)<br>
  :computer: [Rodando o Projeto](#computer-rodando-o-projeto)<br>
  :newspaper: [Gerando documentação com swagger](#newspaper-gerando-documentação-com-swagger)<br>
  :hammer: [Ferramentas](#hammer-ferramentas)<br>
  :clap: [Boas Práticas](#clap-boas-práticas)<br>
  :1234: [Versões](#1234-versões)

<!--te-->
---
<a id="sobre"></a>
## :green_book: Sobre
Projeto de estudo baseado na trilha [Go e Gin: criando API rest com simplicidade](https://www.alura.com.br/curso-online-go-gin-api-rest-simplicidade). Esse projeto tem finalidade puramente didática.


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

Basta rodar o docker-compose:
```bash
$ docker compose up
```
"Aguarde até que as imagens sejam criadas e acesse:

Acesse para API: `http://localhost:8080/alunos`
Acesse para documentação Swagger: `http://localhost:8080/swagger/index.html`


[:arrow_heading_up: voltar](#indice)

---
### :newspaper: Gerando documentação com swagger
Para os desenvolvedores que irão manipular o código ou se inspirar para seus próprios desenvolvimentos, há uma particularidade na documentação Swagger. O comando padrão do [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger) (uma ferramenta que gera documentação Swagger para Go) não consegue ler `structs` que utilizam `gorm.Model`, e isso não está explicitamente mencionado em sua documentação. Pesquisando por uma solução, [encontrei o comando apropriado](https://github.com/swaggo/swag/issues/810) para a geração, que segue abaixo:

```bash
$ swag init --parseDependency --parseInternal
```


<!-- 
https://medium.com/@22vinikaanthwal/generating-swagger-docs-from-go-source-code-2ddc8add251a
https://adityarama1210.medium.com/create-golang-api-doc-with-swag-d73be1767d39
-->
---

<a id="ferramentas"></a>
## :hammer: Ferramentas
As seguintes ferramentas foram usadas na construção do projeto:

- [GVM v1.0.22](https://github.com/moovweb/gvm)
- [Go v1.21.1](https://go.dev/)
- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/index.html)
- [Gin-Swagger](https://github.com/swaggo/gin-swagger)
- [Postgres](https://www.postgresql.org/)
- [Docker 24.0.6](https://www.docker.com/)
- [Docker compose v2.21.0](https://www.docker.com/)
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



