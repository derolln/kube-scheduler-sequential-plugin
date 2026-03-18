.PHONY: build docker-build deploy clean test

IMAGE_NAME ?= sequential-scheduler
IMAGE_TAG ?= latest

build:
	go build -o scheduler .

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

deploy:
	kubectl apply -f manifests/rbac.yaml
	kubectl apply -f manifests/scheduler-config.yaml
	kubectl apply -f manifests/deployment.yaml

undeploy:
	kubectl delete -f manifests/deployment.yaml --ignore-not-found
	kubectl delete -f manifests/scheduler-config.yaml --ignore-not-found
	kubectl delete -f manifests/rbac.yaml --ignore-not-found

clean:
	rm -f scheduler
	go clean

test:
	go test -v ./...

fmt:
	go fmt ./...

vet:
	go vet ./...
