cluster:
	kind create cluster --config=cluster.yaml

build:
	docker build -t device-plugin-1:0 .
	kind load docker-image device-plugin-1:0

start:
	kubectl apply -f plugin.yaml

stop:
	kubectl delete -f plugin.yaml

restart:
	kubectl delete -f plugin.yaml
	kubectl apply -f plugin.yaml
