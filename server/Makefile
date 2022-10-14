TAG=ilingu/tenidditer-api

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker run --env-file=prod.env -it --rm -p 3000:3000 ${TAG}
	@echo Docker images started!

## bup: rebuild all app and starts all containers in the background
bup: build_api up

## build_api: builds the api docker container
build_api:
	@echo Building api container...
	docker build -f tedinitter-api.dockerfile -t ${TAG} .
	@echo Done!

## push: builds the api container and pushs it to dockerHub
push: build_api
	@echo Pushing server binary to DockerHub...
	docker push ${TAG}
	@echo Done!