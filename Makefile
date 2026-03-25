.PHONY: build docker-build deploy clean test

IMAGE_NAME ?= sequential-scheduler
IMAGE_TAG ?= v6
#openssl rand -hex 3 => To generate random tag

build:
	go build -o scheduler .

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

deploy:
	minikube image load sequential-scheduler:$(IMAGE_TAG)
	sed "s/VERSION/$(IMAGE_TAG)/g" manifests/deployment.yaml.tpl > manifests/deployment.yaml
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

logs:
	kubectl -n kube-system logs $$(kubectl get -n kube-system pod -l component=sequential-scheduler -o jsonpath='{.items[0].metadata.name}')
