ifeq (,$(wildcard .env))
    $(error file .env not exsist.)
endif

include .env
export

dev-build: 
	docker build -t $(IMAGE_DEV) -f Dockerfile.dev .

dev-run: dev-build
	docker-compose build && docker-compose up

prod-build: 
	docker build -t $(IMAGE_PROD) -f Dockerfile.prod .

prod-run: prod-build
	docker run -p 8080:8080 -e WEATHER_KEY=$(WEATHER_KEY) -it --rm $(IMAGE_PROD)

test: dev-build
	docker run -it --rm  $(IMAGE_DEV) go test -race ./...

cover: dev-build
	docker run -it --rm $(IMAGE_DEV) go test -cover ./...

deploy:
	chmod +x ./deploy.sh
	./deploy.sh