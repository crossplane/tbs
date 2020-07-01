# Guide

1. Install Crossplane

```
kubectl create namespace crossplane-system

helm repo add crossplane-alpha https://charts.crossplane.io/alpha

helm install crossplane --namespace crossplane-system crossplane-alpha/crossplane
```

2. Install `provider-gcp` and `provider-aws`

```
kubectl apply -f provider-gcp.yaml
kubectl apply -f provider-aws.yaml
```

3. Create `ClusterRole`s

```
kubectl apply -f roles/
```

4. Create `InfrastructureDefinition`s

```
kubectl apply -f definitions/
```

5. Create `Composition`s

```
kubectl apply -f compositions/
```

6. Create EKS Cluster

```
kubectl apply -f cluster-1.yaml
```

6. Create GKE Cluster

```
kubectl apply -f cluster-2.yaml
```

7. Get `kubeconfig`

kubectl get secret gke-conn -o=jsonpath={.data.kubeconfig} | base64 --decode > gke.kube
kubectl get secret eks-conn -o=jsonpath={.data.kubeconfig} | base64 --decode > eks.kube

8. Use `kubeconfig`

kubectl get nodes --kubeconfig=gke.kube
kubectl get nodes --kubeconfig=eks.kube