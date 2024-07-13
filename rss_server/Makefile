SHELL=cmd.exe
NEWS_BINARY=newsApp
API_BINARY=apiApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_news build_api
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"
	@echo "Waiting for 7 seconds..."
	timeout /T 7 /NOBREAK
	make mig_up

## build_news: builds the news-service binary as a linux executable
build_news:
	@echo "Building news binary..."
	chdir ..\news-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${NEWS_BINARY} ./app/cmd
	@echo "Done!"

## build_api: builds the api-service binary as a linux executable
build_api:
	@echo "Building api binary..."
	chdir ..\api-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${API_BINARY} ./app/cmd
	@echo "Done!"

## mig_up: upping database
mig_up:
	goose -dir ./../migrations postgres "user=postgres password=password dbname=rss host=localhost port=5040 sslmode=disable" up


## mig_down: downing database
mig_down:
	goose -dir ./../migrations postgres "user=postgres password=password dbname=rss host=localhost port=5040 sslmode=disable" down