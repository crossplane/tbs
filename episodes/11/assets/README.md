# Guide

1. Install Crossplane and stack-aws
```
kubectl create namespace crossplane-system
helm repo add crossplane-alpha https://charts.crossplane.io/alpha
helm install crossplane --namespace crossplane-system crossplane-alpha/crossplane --set clusterStacks.aws.deploy=true --set clusterStacks.aws.version=v0.6.0
```
2. Create Crossplane AWS Provider
```
./provide.sh
```
3. Create AWS Network Components
```
kubectl apply -f network/ -R
```
4. Create Crossplane infrastructure classes
```
kubectl apply -f infra/ -R
```
5. Install Cert Manager
```
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v0.11.0/cert-manager.yaml
```
6. Install Cluster API
```
kustomize build github.com/kubernetes-sigs/cluster-api//config?ref=v0.3.0-rc.2 | kubectl apply -f -
```
7. Install Kubeadm Bootstrap Provider
```
kustomize build github.com/kubernetes-sigs/cluster-api//bootstrap/kubeadm/config?ref=v0.3.0-rc.2 | kubectl apply -f -
```
8. Install Kubeadm Control Plane Provider
```
kustomize build github.com/kubernetes-sigs/cluster-api//controlplane/kubeadm/config?ref=v0.3.0-rc.2 | kubectl apply -f -
```
9. Install AWS Provider
```
./cluster-api/aws.sh
```
10. Create cluster and control plane
```
kubectl apply -f cluster-api/cluster.yaml
```
11. Connect to control plane
```
kubectl --namespace=default get secret/capi-quickstart-kubeconfig -o json \
  | jq -r .data.value \
  | base64 --decode \
  > ./capi-quickstart.kubeconfig
```
12. Deploy Calico
```
kubectl --kubeconfig=./capi-quickstart.kubeconfig \
  apply -f https://docs.projectcalico.org/v3.12/manifests/calico.yaml
```
13. Create machine node
```
kubectl apply -f cluster-api/machine.yaml
```
14. Encode CAPI kubeconfig in Kubernetes Secret
```
export ENCODE=`base64 capi-quickstart.kubeconfig -w 0`

cat > secret.yaml <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: capi-crossplane
type: Opaque
data:
  kubeconfig: $ENCODE
EOF

kubectl apply -f secret.yaml
```
15. Create a KubernetesTarget that references the Secret
```
cat > target.yaml <<EOF
apiVersion: workload.crossplane.io/v1alpha1
kind: KubernetesTarget
metadata:
  name: capi-crossplane
  labels:
    tbs: eleven
spec:
  connectionSecretRef:
    name: capi-crossplane
EOF

kubectl apply -f target.yaml
```
16. Create MySQL claim
```
kubectl apply -f app/mysql.yaml
```
17. Create App
```
kubectl apply -f app/mysql.yaml
```
18. Interact with App
```bash
curl -d '{"title":"Crossplane+ClusterAPI", "host":"dan and jason", "viewers":1000000}' -H "Content-Type: application/json" -X POST http://<instert-svc-hostname>/create
```
