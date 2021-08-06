Local deployment with minikube
=======

**Description:** Optimization service can be deployed locally in a Kubernetes cluster with minikube for Linux or Docker Desktop for MacOS

![Local K8s deployment](/assets/deployK8s.png)

### Prerequisites

1. Install **minikube** or **Docker Desktop**
2. Install **kubectl**

### Secrets and configuration

Edit environment variable in K8s deployment files to make changes to the configuration of the service. By default it is configured to use local bucket, provided by K8s Persistent Volume Claim.

In order to work with S3 bucket one must provide valid AWS credentials with K8s secrets.

```sh
kubectl create secret generic awssecrets --from-literal=AWS_ACCESS_KEY_ID=*YOUR_AWS_ID* --from-literal=AWS_SECRET_ACCESS_KEY=*YOUR_AWS_SECRET*
```

If pods can't start one of the reasons may be resource shortage. Container limits that can be found in deployment files in k8s folder may be edited or deleted. The number of replicas can be set in configuration files or edited with kubectl commands.

### Build and run

In case of minikube deploy for Linux run:

```sh
make minikube_deploy
```

or take the necessary steps manually:

1. Start minikube `minikube start --cpus 4 --memory 8192`
2. Run `minikube addons enable ingress` to enable ingress plugin for local setup.
3. Run `eval $(minikube -p minikube docker-env)` to set environment variables for current terminal session to point to the minikube's docker images repository. Use `docker images` and `docker ps` commands to check with which repository you are working.
4. Run `docker-compose --env-file=k8s/local.env build` to create local images. Image names and tags are provided in the local.env file. They must match the names in the k8s deployment files.
5. Run `kubectl apply -f k8s` to start the cluster

In case of Docker Desktop for MacOS:

1. Start K8s cluster from Docker Desktop interface.
2. Run `kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.47.0/deploy/static/provider/cloud/deploy.yaml` to enable ingress plugin for local setup. Check if it is running: `kubectl get pods -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx --watch`
3. Run `docker-compose --env-file=k8s/local.env build` to create local images. Image names and tags are provided in the local.env file. They must match the names in the k8s deployment files.
4. Run `kubectl apply -f k8s` to start the cluster.

### Run with helm

1. Install helm using [installation manual](https://helm.sh/docs/intro/install/)

2. Create local images using steps 1-4 or 1-3 respectively from **Build and run** section

3. Run `kubectl apply -f k8s/ingress-service.yaml` to configure ingress service

4. Run `helm install optimization helm/optimization-service/`

You can override default settings with file:

```sh
helm upgrade -f helm/overrride.yaml optimization helm/optimization-service/
```

or set parameter:

```sh
helm upgrade --set api_service.replicas=1 optimization helm/optimization-service/
```

### Test

Watch pods, services and deployments:

```sh
kubectl get pods
kubectl get services
kubectl get deployments
```

Get ip with:

```sh
minikube ip
```

and visit it in browser. It should show the default ui interface of the application.

In case of Docker Desktop just visit `http://localhost`.

### Monitoring

Run

```sh
minikube dashboard
```

to get information about running services, deployments, storages, etc.