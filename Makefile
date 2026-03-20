.PHONY: build docker-build deploy clean test

IMAGE_NAME ?= sequential-scheduler
IMAGE_TAG ?= v2

build:
	go build -o scheduler .

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

deploy:
	kubectl apply -f manifests/rbac.yaml
	kubectl apply -f manifests/scheduler-config.yaml
	kubectl apply -f manifests/deployment.yaml
	minikube image load sequential-scheduler:$(IMAGE_TAG)

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
