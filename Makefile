ifneq ("$(wildcard .env)","")
  include .env
endif

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_BUILD:=$(CURDIR)/build
PROJECT_PATH := $(shell pwd)
GO = $(shell which go)

first-init: install-go-deps

init-env:
	cp $(CURDIR)/.env.example $(CURDIR)/.env

run:
	docker-compose up app

install-go-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/vektra/mockery/v2@v2.43.2
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@v1.16.3

build-app:
	cd cmd/app && go build -o ${LOCAL_BUILD}/app

local-start:
	go run cmd/app/main.go

prepare-tests: ## подготовить окружение к тестированию
	@printf $(COLOR) "Generate ..."
	$(GO) generate ./...

test: prepare-tests ## Запустить тесты
	@printf $(COLOR) "Run tests ..."
	$(GO) test -race -cover -short -v -tags mock -count=2 \
				-coverprofile profile.cov.tmp -p 1 \
				./...
	cat profile.cov.tmp | grep -Ev "_gen.go|mock_|mocks" > profile.cov
	$(MAKE) cover

local-test: prepare-tests ## Запустить все тесты, включая те, что не используют моки
	@printf $(COLOR) "Run tests ..."
	$(GO) test -race -cover -short -v -tags mock,integration -count=1 \
				-coverprofile profile.cov.tmp -p 1 \
				./...
	cat profile.cov.tmp | grep -Ev "_gen.go|mock_|mocks" > profile.cov
	$(MAKE) cover

cover: ## Посчитать coverage проекта
	@printf $(COLOR) "Code coverage ..."
	$(GO) tool cover -func profile.cov

swag:
	./bin/swag init -d internal/server,internal -g router.go -o ./api/swagger --ot yaml

lint:
	./bin/golangci-lint run
