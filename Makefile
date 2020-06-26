run:
	go run main.go

docker-build:
	docker build -t zsiec/duck:$(TAG) .

docker-run:
	docker run -p 8080:8080 zsiec/duck:$(TAG)

docker-push:
	docker push zsiec/duck:$(TAG)

docker-deploy: docker-build docker-push

deploy-dev:
	(cd deploy && ./deploy.sh dev)

delete-dev:
	(cd deploy && ./delete.sh dev)


.PHONY: docker-build docker-run docker-push docker-deploy deploy-dev delete-dev
