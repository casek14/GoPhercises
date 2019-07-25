BIN_NAME = main
CONTAINER_NAME = quay.io/casek14/urlshortenerapp
OLD_MAJOR_VERSION := $(shell cat version.txt | cut -d '.' -f1)
OLD_MINOR_VERSION := $(shell cat version.txt | cut -d '.' -f2)
NEW_TAG := $(shell echo $(OLD_MINOR_VERSION) + 1 | bc)
NEW_CONTAINER := $(CONTAINER_NAME):$(OLD_MAJOR_VERSION).$(NEW_TAG)

.PHONY: help

help:
	@echo "make all - build binary, container and upload it to registry"
	@echo "make build - build golang binary"
	@echo "make build-container - build container with app binary"
	@echo "make upload-container - upload container with latest tag with newer version"
	@echo "make clean - cleanup"

all: build-container upload-container

build: clean
	CGO_ENABLED=0 GOOS=linux 
	go build -a -installsuffix cgo -o $(BIN_NAME) .

build-container: build
	docker build -t $(CONTAINER_NAME):latest .

upload-container: 
	@echo $(CONTAINER_NAME):latest $(NEW_CONTAINER)
	docker tag $(CONTAINER_NAME):latest $(NEW_CONTAINER)
	$(shell echo $(OLD_MAJOR_VERSION).$(NEW_TAG) > version.txt)
	docker push $(NEW_CONTAINER)

clean:
	-rm -f $(BIN_NAME)
