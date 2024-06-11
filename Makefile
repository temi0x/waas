# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=waas
BINARY_PATH=./cmd/waas/main.go

# Docker parameters
DOCKER_IMAGE_NAME=waas
DOCKER_CONTAINER_NAME=waas-container

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(BINARY_PATH)

run:
	./$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test:
	$(GOTEST) ./...

deps:
	$(GOGET) ./...

docker-build:
	docker build -t $(DOCKER_IMAGE_NAME) .

docker-run:
	docker run --name $(DOCKER_CONTAINER_NAME) -d $(DOCKER_IMAGE_NAME)

docker-stop:
	docker stop $(DOCKER_CONTAINER_NAME)

docker-rm:
	docker rm $(DOCKER_CONTAINER_NAME)
