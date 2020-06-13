GO_OS := linux
GO_ARCH := amd64
NETWORK_NAME := local-network
AWS_REGION := ap-northeast-1
EVENT_NAME := scheduled-event
EVENT_DIR := event
ENDPOINT := http://localhost:8000
TABLE_SCHEMA_URI := file://schema/schema.json

.PHONY: all
all: clean build

.PHONY: deps
deps:
	go get -u ./...

.PHONY: clean
clean:
	rm -rf ./bin/*
	go clean

.PHONY: build
build:
	GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build -o ./bin/collect_livers ./functions/collect_livers

.PHONY: test
test: build
	go test -v ./...

create-network:
	docker network inspect $(NETWORK_NAME) > /dev/null 2>&1 || \
	docker network create $(NETWORK_NAME)

create-table:
	@docker-compose up -d
	-aws dynamodb create-table --endpoint-url $(ENDPOINT) --cli-input-json $(TABLE_SCHEMA_URI)

generate-event:
	ls $(EVENT_DIR)/$(EVENT_NAME).json > /dev/null 2>&1 || \
	sam local generate-event cloudwatch $(EVENT_NAME) --region $(AWS_REGION) > $(EVENT_DIR)/$(EVENT_NAME).json

.PHONY: local-invoke
local-invoke: create-network generate-event create-table
	sam local invoke $(NAME) \
	--region $(AWS_REGION) \
	--event $(EVENT_DIR)/$(EVENT_NAME).json \
	--docker-network $(NETWORK_NAME)

