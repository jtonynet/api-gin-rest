# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]
### Added

### Fixed

---

## [0.0.9] - 2023-11-02
### Added
- [Paginação](https://articles.wesionary.team/make-pagination-easy-in-golang-using-pagination-middleware-and-gorm-scope-a5f6eb3bebaa) para a rota de `buscarTodosAlunos`. Melhorias na Paginação em próximos ciclos
- Utilizando [decorator pattern](https://www.henrydu.com/2022/01/05/golang-decorator-pattern/) para cachear o worker insereAluno
- Melhorias na legibilidade geral da codebase

## [0.0.8] - 2023-10-27
### Added
- Adicionado Strategy Pattern para cacheClient
- Adicionado [Redis](https://redis.io/) ao `docker-compose.yml`
- Adicionado classes de cliente `Redis`
- Uso adequado de [custom-middleware](https://gin-gonic.com/docs/examples/custom-middleware/) para gerenciar cache e publicação de mensagen no messageBroker a partir das rotas `Get` e `Post` "decoradas"

### Fixed
- Quando `POST_ALUNO_AS_MESSAGE_FEATURE_FLAG_ENABLED` esta habilitada, worker nao executa e Readiness nao faz checagem de conexão no messageBroker

## [0.0.7] - 2023-10-21
### Added
- Refactor visando melhoria na legibilidade do codigo
- Implementada Reconexão no Message Broker AutoReconnect

## [0.0.6] - 2023-10-17
### Added
- Adicionado `UUID` por aluno
- Adicionado [RabbitMQ](https://www.rabbitmq.com/) ao docker-compose seguindo o [artigo](https://x-team.com/blog/set-up-rabbitmq-with-docker-compose/)
- Criação do `Worker` de processamento de alunos por mensageria
- Feature Flag `POST_ALUNO_AS_MESSAGE_FEATURE_FLAG_ENABLED` para criar aluno de maneira assincrona
- Volumes de `Postgres` e `RabbitMQ` movidos para o diretório `docker_conf`
- `Strategy Pattern` para resolver multiplos message broker

### Fixed
- Adicionada a biblioteca [Exponential Backoff](https://github.com/cenkalti/backoff) para corrigir um bug no ambiente. O `RabbitMQ` demora mais do que o esperado para responder às requisições. Ele informa ao `Docker` que está pronto, mas na verdade, não está, o que ocasiona uma interrupção na API, sem conectividade com esse recurso de infraestrutura. A biblioteca `Backoff` fica responsável por gerenciar as tentativas de conexão pelo período máximo definido na variável de ambiente `API_RETRY_MAX_ELAPSED_TIME_IN_MS`


## [0.0.5] - 2023-10-10
### Added

- Utilizando [gin-contrib/pprof](https://github.com/gin-contrib/pprof) para performar '[Profiling gin with pprof](https://dizzy.zone/2018/08/23/Profiling-gin-with-pprof/)'
- Feature Flag `PPROF_CPU_FEATURE_FLAG_ENABLED` para habilitar o profilling
- Acertos das imagens `Dockerfile`, utilizando `Alpine`
- Alteração na estrutura da `API`, diretório `cmd/api`
- Melhorias no `Readme`.
- Geração  de documentação `Swagger` por comando na imagem agora.

## [0.0.4] - 2023-10-02
### Added

- Utilizando [Gatling v3.9.5](https://gatling.io/) para performar [teste de carga](https://en.wikipedia.org/wiki/Load_testing)
- Acertos das imagens `Dockerfile`
- Geração automática de documentação `Swagger` assim que o `Docker Compose` inicia o projeto.

## [0.0.3] - 2023-09-24
### Added

- Utilizando [Viper](https://github.com/spf13/viper) para a gestão de variáveis de ambiente
- Diagrama `Mermaid` simples adicionado ao `readme`
- Rotas `http://localhost:8080/readiness` e `http://localhost:8080/liveness`
 
## [0.0.2] - 2023-09-23
### Added

- Utilizando [gin-swagger](https://github.com/swaggo/gin-swagger) para documentar a `API` no padrão `OpenAPI`

## [0.0.1] - 2023-09-23
### Added

- Projeto com base no curso [Go e Gin: criando API rest com simplicidade](https://www.alura.com.br/curso-online-go-gin-api-rest-simplicidade) finalizado de acordo com a trilha seguida
- API Gin `dockerizada`

[0.0.9]: https://github.com/jtonynet/api-gin-rest/compare/v0.0.8...v0.0.9
[0.0.8]: https://github.com/jtonynet/api-gin-rest/compare/v0.0.7...v0.0.8
[0.0.7]: https://github.com/jtonynet/api-gin-rest/compare/v0.0.6...v0.0.7
[0.0.6]: https://github.com/jtonynet/api-gin-rest/compare/v0.0.5...v0.0.6
[0.0.5]: https://github.com/jtonynet/api-gin-rest/compare/v0.0.4...v0.0.5
[0.0.4]: https://github.com/jtonynet/api-gin-rest/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/jtonynet/api-gin-rest/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/jtonynet/api-gin-rest/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/jtonynet/api-gin-rest/releases/tag/v0.0.1
