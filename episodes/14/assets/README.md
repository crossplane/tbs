# Guide

1. Install Crossplane with provider-gcp

```
kubectl create namespace crossplane-system
helm repo add crossplane-master https://charts.crossplane.io/master/
version=$(helm search repo crossplane --devel | awk '$1 == "crossplane-master/crossplane" {print $2}')
helm install crossplane --namespace crossplane-system crossplane-master/crossplane --version $version --set clusterStacks.gcp.deploy=true --set clusterStacks.gcp.version=master --devel
```

2. Create GCP `Provider`

Documentation
[here](https://crossplane.io/docs/v0.8/cloud-providers/gcp/gcp-provider.html).

## Plain OPA

1. Setup OPA Configuration

Generate certificate authority and key pair:
```
kubectl create namespace opa

openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -days 100000 -out ca.crt -subj "/CN=admission_ca"
```

Generate TLS config:
```
cat >server.conf <<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, serverAuth
EOF
```

```
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj "/CN=opa.opa.svc" -config server.conf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 100000 -extensions v3_req -extfile server.conf
```

Create TLS secret:
```
kubectl -n opa create secret tls opa-server --cert=server.crt --key=server.key
```

2. Deploy OPA

Deploy the admission controller:
```
kubectl -n opa apply -f plain/admission-controller.yaml
```

Register as validating webhook:
```
cat > plain/webhook-configuration.yaml <<EOF
kind: ValidatingWebhookConfiguration
apiVersion: admissionregistration.k8s.io/v1beta1
metadata:
  name: opa-validating-webhook
webhooks:
  - name: validating-webhook.openpolicyagent.org
    namespaceSelector:
      matchExpressions:
      - key: openpolicyagent.org/webhook
        operator: NotIn
        values:
        - ignore
    rules:
      - operations: ["CREATE", "UPDATE"]
        apiGroups: ["*"]
        apiVersions: ["*"]
        resources: ["*"]
    clientConfig:
      caBundle: $(cat ca.crt | base64 | tr -d '\n')
      service:
        namespace: opa
        name: opa
EOF
```

```
kubectl apply -f plain/webhook-configuration.yaml
```

Ignore the `opa` and `kube-system` namespaces:
```
kubectl label ns kube-system openpolicyagent.org/webhook=ignore
kubectl label ns opa openpolicyagent.org/webhook=ignore
```

3. Define Policy

```
kubectl -n opa create configmap db-size --from-file=./plain/db-size.rego
```

### Static Provisioning

1. Create invalid `CloudSQLInstance`

```
kubectl apply -f plain/cloudsql-invalid.yaml
```

2. Create valid `CloudSQLInstance`

```
kubectl apply -f plain/cloudsql-valid.yaml
```

### Dynamic Provisioning

1. Create `CloudSQLInstanceClass`

```
kubectl apply -f plain/cloudsqlclass-invalid.yaml
```

2. Create `MySQLInstance` claim

```
kubectl apply -f plain/mysqlclaim.yaml
```

3. Delete DB Size policy

```
kubectl -n opa delete configmap db-size
```

## Gatekeeper

1. Deploy Gatekeeper

```
kubectl apply -f https://raw.githubusercontent.com/open-policy-agent/gatekeeper/master/deploy/gatekeeper.yaml
```

2. Create `ConstraintTemplate` for global CloudSQL policy

```
kubectl apply -f gatekeeper/constraint-template.yaml
```

3. Create `Constraint` for global CloudSQL policy

```
kubectl apply -f gatekeeper/constraint.yaml
```

4. Create `ConstraintTemplate` for namespace-level policy

```
kubectl apply -f gatekeeper/constraint-template-app.yaml
```

4. Create `Constraints` for namespace-level policy

```
kubectl apply -f gatekeeper/constraint-app.yaml
```

### Static Provisioning

1. Create invalid `CloudSQLInstance`

```
kubectl apply -f plain/cloudsql-invalid.yaml
```

2. Create valid `CloudSQLInstance`

```
kubectl apply -f plain/cloudsql-valid.yaml
```

### Dynamic Provisioning

1. Create `CloudSQLInstanceClass`

```
kubectl apply -f plain/cloudsqlclass-invalid.yaml
```

2. Create `MySQLInstance` claim

```
kubectl apply -f plain/mysqlclaim.yaml
```

## OPA in Remote Cluster

1. Create `GKEClusterClass`

```
kubectl apply -f remote/gkeclusterclass.yaml
```

2. Create `KubernetesCluster` claim

```
kubectl apply -f remote/k8scluster.yaml
```

3. Deploy Gatekeeper into remote cluster

```
curl https://raw.githubusercontent.com/open-policy-agent/gatekeeper/master/deploy/gatekeeper.yaml | kubectl crossplane pack - | kubectl apply -f -
```

4. Deploy `ConstraintTemplate` into remote cluster

```
cat remote/constraint-template.yaml | kubectl crossplane pack - | kubectl apply -f -
```

5. Deploy `Constraint` into remote cluster

```
cat remote/constraint.yaml | kubectl crossplane pack - | kubectl apply -f -
```

6. Create invalid `CloudSQLInstance` in remote cluster

```
cat remote/pod.yaml | kubectl crossplane pack - | kubectl apply -f -
```