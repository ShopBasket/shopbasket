.PHONY: build
build:
	go build

.PHONY: docker-image
docker-image:
	docker build -t ghcr.io/shopbasket/shopbasket:latest .

.PHONY: docker-image-push
docker-image-push:
	docker push ghcr.io/shopbasket/shopbasket:latest
