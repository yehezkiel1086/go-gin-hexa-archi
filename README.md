# Go Gin Hexa Architecture Template

Golang Gin Hexagonal Architecture Implementation of employees repository of a company. 

Tech stack:
- Language: Go
- Framework: Gin
- Architecture: Hexagonal (clean) architecture
- DB: Postgres
- DBMS: GORM
- Caching: Redis
- Containerization: Docker
- CI/CD: Github Actions 

Roles:
- Employee (2001)
- Admin (5150)

## Requirements

Go >= 1.24 \
Docker

## Installation

1. Run the docker containers
    ```sh
    docker compose up -d
    ```
2. Download the go libraries
    ```sh
    go mod tidy
    ```
3. Run the api server
    ```sh
    go run cmd/main.go
    ```

or you could also utilize the existing `Makefile`
