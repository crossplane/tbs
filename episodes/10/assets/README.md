# Guide

This guide goes through setting up two different Kubernetes local development environments using Crossplane and Okteto!

## Steps

### 1. Install Okteto CLI

If you are running MacOS / Linux on your local machine, the Okteto CLI can be installed with a single command:

```bash
curl https://get.okteto.com -sSfL | sh
```

Otherwise, take a look at the Okteto [docs](https://okteto.com/docs/getting-started/installation/index.html) to install on your OS of choice.

### 2. Install Crossplane & stack-gcp

For full Crossplane installation instruction, take a look at the [docs](https://crossplane.io/docs/master/install-crossplane.html). Installing Crossplane with stack-gcp using Helm 3 is described below:

```bash
kubectl create namespace crossplane-system
helm repo add crossplane-master https://charts.crossplane.io/master/
helm search repo crossplane --devel # get version
helm install crossplane --namespace crossplane-system crossplane-master/crossplane --version <version> --set clusterStacks.gcp.deploy=true --set clusterStacks.gcp.version=master --devel
```

### 3. Run Simple Golang App

The `/app1` directory includes the same sample code used in the Okteto Golang [Quick Start](https://okteto.com/docs/samples/golang/index.html). Okteto will create the `Deployment` for you, but we do it manually here.

Start by creating a `Namespace` for our local development:

```bash
kubectl create namespace okteto
```

Now create the deployment:

```bash
kubectl apply -f k8s.yaml
```

You should see the `hello-world-*` `Pod` running in the `okteto` `Namespace`:

```bash
kubectl get pods -n okteto
```

Now we want to start out local development environment using Okteto:

```bash
okteto up -n okteto
```

Once setup completes, you should be in a new bash shell. If you run `ls` you should see the contents of `/app1` in your current directory. We can now run our application:

```bash
okteto> go run main.go
```

You should be able to navigate to `localhost:8080` and be greeted by your web server message!

When finished, exit the shell session and run `okteto down`. Then clean up our deployment:

```bash
kubectl delete -f deployment.yaml
```

### 4. Run Golang App with MySQL Database

Now we want to use a slightly more complex application that talks to an external MySQL database. All of the code and manifests for this scenario are in `/app2`.

First we need to create our `Provider` object and its `Secret` so we can communicate with GCP. Either replace `BASE64ENCODED_GCP_PROVIDER_CREDS` and `PROJECT_ID` in `provider.yaml` or take a look at the stack-gcp [docs](https://crossplane.io/docs/master/cloud-providers/gcp/gcp-provider.html) for more information.

Now we can dynamically provisioning our database using Crossplane:

```bash
kubectl apply -f cloudsqlclass.yaml
kubectl apply -f mysql.yaml
```

When the database is ready you should see that the `Status` of our `MySQLInstance` is `Bound`. You should also see that a connection `Secret` name `mysqlconn` has been created in the `okteto` `Namespace`.

```bash
kubectl get secrets -n okteto
```

We inject the data in this `Secret` into our `Deployment` to be able communicate with the database. We are going to once again start off our local development process by creating our `Deployment`.

```bash
kubectl apply -f k8s.yaml
```

You can check that the pod starts successfully again. This one may take a bit more time to reach status `Running` due to the fact that it uses an `InitContainer` to create the `okteto` database on the Cloud SQL instance.

```bash
kubectl get pods -n okteto
```

We can once again start out local development environment and start our application:

```bash
okteto up -n okteto
```

```bash
okteto> go run main.go
```

If you navigate to `localhost:8080` you should be greeted with `Welcome to the 11th episode of The Binding Status!`.

Now let's see if we are connecting to the database as expected. If you navigate to `localhost:8080/list` you should get a response of `[]`. This makes sense as we haven't written any records to the database. Let's change that:

```bash
curl -d '{"title":"Crossplane+Okteto", "host":"dan and ramiro", "viewers":1000000}' -H "Content-Type: application/json" -X POST http://localhost:8080/create
```

Now if we navigate to `localhost:8080/list` we should get an array with our one record in it!

### 5. Clean Up

We can clean up as we did earlier by running `okteto down` and `kubectl delete -f k8s.yaml`. We also want to clean up our Cloud SQL instance. Because our `MySQLInstance` claim is bound to a `CloudSQLInstance` with `Reclaim: Delete`, the resource will be cleaned up when the claim is deleted:

```bash
kubectl delete -f mysql.yaml
```