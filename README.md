# Estudo API Rest em Golang com Gin
<br> 
<img src="./docs/images/gopher_2.png">

[<img src="./docs/images/icons/go.svg" width="25px" height="25px" alt="go" title="Go"> <img src="./docs/images/icons/docker.svg" width="25px" height="25px" alt="Docker" title="Docker"> <img src="./docs/images/icons/github.svg" width="25px" height="25px" alt="GitHub" title="GitHub"> <img src="./docs/images/icons/visualstudiocode.svg" width="25px" height="25px" alt="vscode" title="vscode">](#estudo-de-autenticação-testes-e-segurança-em-nodejs) <!-- icons by https://simpleicons.org/?q=types -->



![Badge Status](https://img.shields.io/badge/STATUS-EM_DESENVOLVIMENTO-green)

---

<a id="indice"></a>
## :arrow_heading_up: Índice
<!--ts-->
- [Go: Go e Gin: criando API rest com simplicidade](#estudo-api-rest-em-golang-com-gin)<br>
  :arrow_heading_up: [Índice](#arrow_heading_up-índice)<br>
  :green_book: [Sobre](#green_book-sobre)<br>
  :computer: [Rodando o Projeto](#computer-rodando-o-projeto)<br>
  :hammer: [Ferramentas](#hammer-ferramentas)<br>
  :clap: [Boas Práticas](#clap-boas-práticas)<br>

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

Primeiro inicie o docker-compose de infra/bancos:
```bash
$ docker-compose up
```
Compile e rode o binário:

```bash
$ go run main.go
```

Acesse: `http://localhost:8000/api/personalidades`


[:arrow_heading_up: voltar](#indice)

---

<a id="ferramentas"></a>
## :hammer: Ferramentas
As seguintes ferramentas foram usadas na construção do projeto:

- [GVM v1.0.22](https://github.com/moovweb/gvm)
- [Go v1.21.1](https://go.dev/)
- [Docker 24.0.6](https://www.docker.com/)
- [Docker compose 1.29.2](https://www.docker.com/)

[:arrow_heading_up: voltar](#indice)

---


<a id="boas-praticas"></a>
## :clap: Boas Práticas
Seguindo boas práticas de desenvolvimento:
- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)

[:arrow_heading_up: voltar](#indice)



