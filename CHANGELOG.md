# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added


## [0.0.5] - 2023-10-02
### Added

- Utilizando [gin-contrib/pprof](https://github.com/gin-contrib/pprof) para performar [Profiling gin with pprof](https://dizzy.zone/2018/08/23/Profiling-gin-with-pprof/)
- Feature Flag `PPROF_CPU_FEATURE_FLAG_ENABLED` para habilitar o profilling
- Acertos das imagens Dockerfile, utilizando Alpine
- Alteracao na estrutura da API, diretorio `cmd/api`
- Melhorias no Readme.

## [0.0.4] - 2023-10-02
### Added

- Utilizando [Gatling v3.9.5](https://gatling.io/) para performar [teste de carga](https://en.wikipedia.org/wiki/Load_testing)
- Acertos das imagens Dockerfile
- Geração automática de documentação Swagger assim que o Docker Compose inicia o projeto.

## [0.0.3] - 2023-09-24
### Added

- Utilizando [Viper](https://github.com/spf13/viper) para a gestão de variáveis de ambiente
- Diagrama mermaid simples adicionado ao reademe
- Rotas `http://localhost:8080/readiness` e `http://localhost:8080/liveness`
 
## [0.0.2] - 2023-09-23
### Added

- Utilizando [gin-swagger](https://github.com/swaggo/gin-swagger) para documentar a API no padrão OpenAPI

## [0.0.1] - 2023-09-23
### Added

- Projeto com base no curso [Go e Gin: criando API rest com simplicidade](https://www.alura.com.br/curso-online-go-gin-api-rest-simplicidade) finalizado de acordo com a trilha seguida
- API Gin dockerizada

[0.0.5]: https://github.com/jtonynet/api-gin-rest/compare/v0.0.4...v0.0.5
[0.0.4]: https://github.com/jtonynet/api-gin-rest/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/jtonynet/api-gin-rest/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/jtonynet/api-gin-rest/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/jtonynet/api-gin-rest/releases/tag/v0.0.1
