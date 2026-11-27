# Atalhos do projeto. Funciona como os "scripts" do package.json.
#
#   make run CAP=01   ->  roda capitulos/01-*/main.go no Docker
#   make ex  CAP=01   ->  roda os testes dos exercícios do capítulo 01

CAP ?= 00
DIR := $(firstword $(wildcard capitulos/$(CAP)-*))
GO  := docker compose run --rm go

.DEFAULT_GOAL := help
.PHONY: help list run ex test test-all fmt vet shell image check-cap todo todo-test

help: ## Mostra os comandos disponíveis
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

list: ## Lista os capítulos
	@ls -1 capitulos

run: check-cap ## Roda o main.go de um capítulo (ex: make run CAP=01)
	$(GO) go run ./$(DIR)

ex: check-cap ## Roda os testes dos exercícios de um capítulo (ex: make ex CAP=01)
	$(GO) go test -v ./$(DIR)/exercicios/...

test: check-cap ## Roda todos os testes de um capítulo
	$(GO) go test -v ./$(DIR)/...

test-all: ## Roda os testes do repositório inteiro
	$(GO) go test ./...

fmt: ## Formata todo o código com gofmt
	$(GO) gofmt -l -w .

vet: ## Análise estática com go vet
	$(GO) go vet ./...

shell: ## Abre um terminal dentro do container Go
	$(GO) bash

image: ## (Re)constrói a imagem Docker
	docker compose build

todo: ## Sobe a API de tarefas em http://localhost:8080 (Ctrl+C para parar)
	docker compose up --build todo-api

todo-test: ## Roda os testes da API de tarefas
	$(GO) go test -v ./projetos/todo-api/...

check-cap:
	@if [ -z "$(DIR)" ]; then \
		echo "Capítulo '$(CAP)' não encontrado. Use 'make list' para ver os disponíveis."; \
		exit 1; \
	fi
