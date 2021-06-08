OPTIMIZATION_SERVER_NAME = optimization_server
API_SERVER_NAME = api_server

GOARCH = amd64

.PHONY:

help: ## Display this help screen
	@echo "Makefile available targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  * \033[36m%-15s\033[0m %s\n", $$1, $$2}'

dep: ## Download all dependencies
	go mod tidy && go mod vendor
	
build: dep lint build_optimization build_api ## Build all

build_optimization: ## Build optimization service locally
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -mod vendor -o optimization_server/bin/${OPTIMIZATION_SERVER_NAME} ./optimization_server/run

run_optimization: lint ## Run optimization service locally
	cd optimization_server && ./rundev.sh

build_api: ## Build api service
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -mod vendor -o api_server/bin/${API_SERVER_NAME} ./api_server/run

run_api: lint ## Run api service locally
	cd api_server && ./rundev.sh

test: lint ## Run tests
	go test -race -p 1 -timeout 300s -coverprofile=.test_coverage.txt ./... && \
    	go tool cover -func=.test_coverage.txt | tail -n1 | awk '{print "Total test coverage: " $$3}'
	@rm .test_coverage.txt

lint: ## Lint the source files
	golangci-lint run

docker_build: ## Build optimization engine service with docker-compose

docker_run: ## Run optimization engine service with in docker-compose

docker_push: ## Pushes the build images to the AWS ECR

		