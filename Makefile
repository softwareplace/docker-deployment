output ?= ./ci

docker-build:
	sh ./deployment_builder $(output)