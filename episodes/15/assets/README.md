# Guide

## (Optional) Setting up a Kind cluster

Make sure you have installed [kind](https://kind.sigs.k8s.io/)

```
kind create cluster
```

## Installing Knative with Mink

```
kubectl apply -f mink.yaml
```

If using local Kind cluster, install local variant of Mink
```
kubectl apply -f mink-local.yaml
```

## Installing Crossplane

Install Crossplane from `alpha` channel

```
kubectl create namespace crossplane-system
helm repo add crossplane-alpha https://charts.crossplane.io/alpha
helm install crossplane --namespace crossplane-system crossplane-alpha/crossplane --set clusterStacks.gcp.deploy=true --set clusterStacks.gcp.version=v0.9.0
```

## Configuring GCP credentials

Instructions for setting up your GCP credentials can be found [here](https://crossplane.io/docs/v0.10/cloud-providers/gcp/gcp-provider.html).

## Demo 1: Consuming a CloudSQL database from a Knative Service

Create your `CloudSQLClass`
```
kubectl apply -f cloudsqlclass.yaml
```

Create your `PostgreSQLInstance`
```
cd service
kubectl apply -f psql-claim.yaml
```

Build and push image for Knative `Service`
```
docker build . -t hasheddan/tbs15service:latest
docker push hasheddan/tbs15service:latest
```

Deploy Knative `Service`
```
kubectl apply -f service.yaml
```

Get service URL
```
export APP_URL=$(kubectl get ksvc | awk '$1 == "tbs-live" {print $2}')
```

Interact with service
```bash
curl -d '{"name":"Daniel Mangum", "location":"St. Louis, MO", "twitter":"hasheddan"}' -H "Content-Type: application/json" -X POST $APP_URL/create
```

## Demo 2: Automatically injecting connection information using bindings

Modify connection secret
```
export DB_USER=$(kubectl get secret conn -o=jsonpath={.data.username} | base64 --decode -)
export DB_PASSWORD=$(kubectl get secret conn -o=jsonpath={.data.password} | base64 --decode -)
export DB_HOST=$(kubectl get secret conn -o=jsonpath={.data.endpoint} | base64 --decode -)

kubectl create secret generic conn-bind --from-literal=connectionstr=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/postgres
```

Create `SQLBinding`
```
kubectl apply -f binding.yaml
```

Create Knative `Service`
```
kubectl apply -f service-bind.yaml
```
