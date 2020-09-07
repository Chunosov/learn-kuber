#!/bin/sh

echo "Installing knative ..."
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.17.0/serving-crds.yaml
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.17.0/serving-core.yaml

echo "Installing kong ..."
kubectl apply --filename https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/0.9.x/deploy/single/all-in-one-dbless.yaml

echo "Applying kong Ingress Controller ..."
kubectl patch configmap/config-network --namespace knative-serving --type merge --patch '{"data":{"ingress.class":"kong"}}'

echo "Done"
