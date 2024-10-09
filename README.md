# Champi-maker

Este projeto é uma aplicação backend escrita em Go, seguindo a Arquitetura Hexagonal. Ela fornece uma API RESTful para gerenciar usuários, times, campeonatos, partidas e estatísticas.

## Sumário

- [Funcionalidades](#funcionalidades)
- [Arquitetura](#arquitetura)
- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Iniciando o Projeto](#iniciando-o-projeto)
  - [Pré-requisitos](#pré-requisitos)
  - [Instalação](#instalação)
  - [Configuração](#configuração)
- [Uso](#uso)
  - [Executando a Aplicação](#executando-a-aplicação)
  - [Executando Testes](#executando-testes)
  - [Compilando a Aplicação](#compilando-a-aplicação)
  - [Migrações do Banco de Dados](#migrações-do-banco-de-dados)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Contribuição](#contribuição)
- [Licença](#licença)

## Funcionalidades

- Registro e autenticação de usuários com JWT
- Operações CRUD para Times
- Criação e gerenciamento de Campeonatos
- Geração e gerenciamento de Partidas dentro dos Campeonatos
- Atualizações em tempo real via filas de mensagens RabbitMQ
- Testes unitários e de integração abrangentes
- Manipulação segura de senhas e autenticação
- API RESTful seguindo as melhores práticas

## Arquitetura

A aplicação é construída seguindo a Arquitetura Hexagonal, promovendo uma clara separação de responsabilidades e tornando o código mais fácil de manter e testar.

- **Camada de Domínio**: Contém a lógica de negócios e as entidades do domínio.
- **Camada de Aplicação**: Implementa os casos de uso e orquestra o fluxo entre as camadas.
- **Camada de Infraestrutura**: Lida com preocupações externas, como acesso ao banco de dados e mensagens.
- **Camada de Interfaces**: Contém os handlers HTTP, middlewares e outros adaptadores de interface.

## Tecnologias Utilizadas

- **Linguagem**: Go (Golang)
- **Framework Web**: Gin
- **Banco de Dados**: PostgreSQL
- **Driver do Banco**: pgx (sem ORM)
- **Mensageria**: RabbitMQ
- **Autenticação**: JWT
- **Testes**: Pacote `testing`, `testify`
- **Linting**: golangci-lint
- **Ferramenta de Migração**: [migrate](https://github.com/golang-migrate/migrate)

## Iniciando o Projeto

### Pré-requisitos

- Go (versão 1.16 ou superior)
- PostgreSQL
- RabbitMQ
- Git
- make (se for utilizar o Makefile)
- [golang-migrate](https://github.com/golang-migrate/migrate) (para migrações do banco)
- [golangci-lint](https://golangci-lint.run/) (para linting)

### Instalação

1. **Clone o repositório**

```bash
   git clone https://github.com/seu_usuario/seu_projeto.git
   cd seu_projeto
```

2. **Instale as dependências do Go**

```bash
   go mod download
```

### Configuração

Crie um arquivo .env na raiz do seu projeto e defina as seguintes variáveis de ambiente:

```env
    # Configurações do Banco de Dados
    DATABASE_URL_TEST=postgres://usuario:senha@localhost:5432/seu_banco_teste?sslmode=disable
```

- Substitua usuario, senha e seu_banco com suas credenciais reais do banco de dados.

## Uso

### Executando a Aplicação

Você pode executar a aplicação usando o Makefile:

```bash
   make run
```

Ou diretamente com o Go:

```bash
   go run cmd/api/main.go
```

### Executando Testes

Para executar todos os testes:

```bash
   make test
```

### Compilando a Aplicação

Para compilar o executável da aplicação:

```bash
   make build
```

- Isso criará um executável chamado app_name.

### Migrações do Banco de Dados

Certifique-se de ter a ferramenta migrate instalada.

- Aplicar migrações:

```bash
   make migrate-up
```

- Reverter migrações:

```bash
   make migrate-down
```

## Estrutura do Projeto

```plaintext
.
├── cmd
│   └── api
│       └── main.go          # Ponto de entrada da aplicação
├── internal
│   ├── domain
│   │   ├── entity           # Entidades do domínio
│   │   └── repository       # Interfaces dos repositórios
│   ├── application
│   │   └── service          # Lógica de negócios e casos de uso
│   ├── infrastructure
│   │   ├── repository       # Implementações de acesso ao banco
│   │   ├── messaging        # Implementações de mensageria com RabbitMQ
│   │   ├── db               # Implementações do banco de dados
│   │   └── config           # Implementações para ler a .env
│   └── interfaces
│       ├── handler          # Handlers HTTP
│       └── middleware       # Middlewares HTTP (ex: autenticação)
├── go.mod
├── go.sum
├── Makefile
└── README.md

```