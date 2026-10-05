# Sprint 2: Container Orchestration with Kubernetes

## Project Title

**Kubernetes Deployment and Orchestration of a Go Application**

## Description

This project is part of the DevOps Internship 2026 evaluation. The objective of Sprint 2 was to deploy the containerized Go application from Sprint 1 onto a local Kubernetes cluster using Minikube.

The application was deployed using a Kubernetes Deployment with two replicas. A ConfigMap and Secret were injected into the application container as environment variables. Kubernetes liveness and readiness probes were configured using the application's `/health` and `/ready` endpoints.

The application was exposed using a NodePort Service.

## Project Objectives

* Deploy the existing Go application to Kubernetes.
* Run two application replicas using a Deployment.
* Use a ConfigMap for application configuration.
* Use a Secret for sensitive configuration.
* Inject ConfigMap and Secret values as environment variables.
* Configure a liveness probe using `/health`.
* Configure a readiness probe using `/ready`.
* Expose the application using a NodePort Service.
* Verify the Deployment, Pods, Service, endpoints, probes, and application health.

## Project Structure

```text
sprint-2/
├── README.md
└── manifests/
    ├── configmap.yaml
    ├── secret.yaml
    ├── deployment.yaml
    └── service.yaml
```

## Getting Started

### Dependencies

The project was tested using:

* Ubuntu 24.04 WSL
* Docker
* kubectl v1.37.1
* Minikube v1.39.0
* Kubernetes v1.37.0
* Git

The Go application image was already built during Sprint 1, so Go did not need to be installed separately on the host for Sprint 2.

### Installing

The project was located at:

```text
~/Devops-Intership-2026/sprint-2
```

The Sprint 2 application image was:

```text
sprint2-go-app:latest
```

## Executing Program

### Step 1: Verify kubectl

The installed kubectl version was checked with:

```bash
kubectl version --client
```

Result:

```text
Client Version: v1.37.1
Kustomize Version: v5.8.1
```

### Step 2: Verify Minikube

Minikube version was checked with:

```bash
minikube version
```

Minikube was installed as version:

```text
minikube version: v1.39.0
```

### Step 3: Start Minikube

The local Kubernetes cluster was started using the Docker driver:

```bash
minikube start --driver=docker
```

The cluster was verified with:

```bash
minikube status
```

The final cluster status was:

```text
host: Running
kubelet: Running
apiserver: Running
kubeconfig: Configured
```

The Kubernetes node was verified with:

```bash
kubectl get nodes
```

The node was in `Ready` state and running Kubernetes `v1.37.0`.

### Step 4: Verify Kubernetes System Pods

The Kubernetes system pods were checked using:

```bash
kubectl get pods -A
```

The required Kubernetes system components were running successfully.

### Step 5: Verify the Application Docker Image

The Sprint 2 application image was verified using:

```bash
docker images | grep sprint2-go-app
```

The image used for the deployment was:

```text
sprint2-go-app:latest
```

The image size was approximately `6.02 MB`.

### Step 6: Load the Image into Minikube

Because the application image was available locally rather than in a remote registry, it was loaded into Minikube using:

```bash
minikube image load sprint2-go-app:latest
```

The image was verified inside Minikube using:

```bash
minikube image ls | grep -i sprint
```

The output included:

```text
docker.io/library/sprint2-go-app:latest
```

### Step 7: Validate the Kubernetes Manifests

The manifests were validated before deployment using:

```bash
kubectl apply --dry-run=client -f manifests/
```

The validation completed successfully for:

```text
configmap/sprint2-config created (dry run)
deployment.apps/sprint2-go-deployment created (dry run)
secret/sprint2-secret created (dry run)
service/sprint2-go-service created (dry run)
```

### Step 8: Deploy the Kubernetes Resources

The Kubernetes resources were created using:

```bash
kubectl apply -f manifests/
```

The resources created were:

```text
configmap/sprint2-config created
deployment.apps/sprint2-go-deployment created
secret/sprint2-secret created
service/sprint2-go-service created
```

### Step 9: Verify the Deployment

The Deployment was checked using:

```bash
kubectl get deployments
```

The final Deployment state was:

```text
NAME                    READY   UP-TO-DATE   AVAILABLE
sprint2-go-deployment   2/2     2            2
```

This confirms that the requested two replicas were running and available.

### Step 10: Verify the Pods

The application pods were checked using:

```bash
kubectl get pods
```

The two running pods were:

```text
sprint2-go-deployment-5b76b8454c-fl9wv   1/1   Running   0
sprint2-go-deployment-5b76b8454c-jqdz6   1/1   Running   0
```

Both pods were ready and had zero restarts.

Pod status was also monitored using:

```bash
kubectl get pods -w
```

Both pods reached `1/1 Running` successfully.

### Step 11: Verify the Service

The Kubernetes Services were checked using:

```bash
kubectl get services
```

The application Service was:

```text
sprint2-go-service   NodePort   10.98.45.219   <none>   3000:30080/TCP
```

The Service configuration was:

```text
Service type: NodePort
Service port: 3000
Target port: 3000
NodePort: 30080
```

### Step 12: Verify Service Endpoints

The application Service endpoints were checked using:

```bash
kubectl get endpoints sprint2-go-service
```

The endpoints were:

```text
10.244.0.3:3000,10.244.0.4:3000
```

This confirmed that the Service had both application pods as backend endpoints.

### Step 13: Get the Application URL

The application was exposed through Minikube using:

```bash
minikube service sprint2-go-service --url
```

During the verification session, Minikube returned:

```text
http://127.0.0.1:44637
```

The URL is generated by Minikube for the active session and can change after restarting the service or Minikube. The command above should therefore be used whenever the current URL is needed.

### Step 14: Test the Application

The application home page was opened using the Minikube service URL.

The application was successfully accessible through the NodePort Service.

### Step 15: Test the Health Endpoint

The health endpoint was tested at:

```text
http://127.0.0.1:44637/health
```

The endpoint returned HTTP `200 OK` and displayed the application's healthy status.

The same `/health` endpoint is used by the Kubernetes liveness probe.

### Step 16: Test the Readiness Endpoint

The readiness endpoint was tested at:

```text
http://127.0.0.1:44637/ready
```

The endpoint returned:

```text
READY
```

The same `/ready` endpoint is used by the Kubernetes readiness probe.

### Step 17: Inspect the Pod Configuration

The pod was inspected using:

```bash
kubectl describe pod sprint2-go-deployment-5b76b8454c-fl9wv
```

The inspection confirmed the following:

```text
Image: sprint2-go-app:latest
Port: 3000/TCP
Liveness: http-get http://:3000/health
Readiness: http-get http://:3000/ready
ConfigMap: sprint2-config
Secret: sprint2-secret
Ready: True
Restart Count: 0
```

The pod was successfully running and ready.

## Kubernetes Configuration

### ConfigMap

The ConfigMap was created in:

```text
manifests/configmap.yaml
```

It contains the following non-sensitive configuration:

```text
APP_ENV=production
APP_NAME=sprint2-go-app
LOG_LEVEL=info
```

The Deployment references the ConfigMap using:

```yaml
envFrom:
  - configMapRef:
      name: sprint2-config
```

The ConfigMap was verified with:

```bash
kubectl get configmap sprint2-config
```

### Secret

The Secret was created in:

```text
manifests/secret.yaml
```

It contains the demonstration application secret.

The Deployment references it using:

```yaml
envFrom:
  - secretRef:
      name: sprint2-secret
```

The Secret was verified with:

```bash
kubectl get secret sprint2-secret
```

The Secret value was not included in screenshots or documentation.

> The Secret in this project is only a demonstration value and is not a real production credential.

### Deployment

The Deployment was created in:

```text
manifests/deployment.yaml
```

Important configuration:

```text
Deployment name: sprint2-go-deployment
Replicas: 2
Container name: sprint2-go-app
Container image: sprint2-go-app:latest
Container port: 3000
Image pull policy: IfNotPresent
```

### Liveness Probe

The liveness probe checks the application's `/health` endpoint:

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 3000
```

The configured values were:

```text
Initial delay: 5 seconds
Period: 10 seconds
Timeout: 2 seconds
Failure threshold: 3
```

### Readiness Probe

The readiness probe checks the application's `/ready` endpoint:

```yaml
readinessProbe:
  httpGet:
    path: /ready
    port: 3000
```

The configured values were:

```text
Initial delay: 5 seconds
Period: 5 seconds
Timeout: 2 seconds
Failure threshold: 3
```

### Container Security Configuration

The Kubernetes Deployment also included:

```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 10001
  allowPrivilegeEscalation: false
```

This matches the non-root user used by the Sprint 1 Docker image.

### Service

The Service was created in:

```text
manifests/service.yaml
```

The actual Service configuration was:

```text
Service name: sprint2-go-service
Type: NodePort
Port: 3000
Target port: 3000
NodePort: 30080
Selector: app=sprint2-go-app
```

## Verification Summary

The following checks were completed successfully:

| Check | Result |
|---|---|
| Minikube | Running |
| Kubernetes node | Ready |
| Application image | Available in Minikube |
| Deployment replicas | 2/2 |
| Application pods | 2/2 Running |
| Pod readiness | Ready |
| Pod restarts | 0 |
| ConfigMap | Created and referenced |
| Secret | Created and referenced |
| Liveness probe | `/health` |
| Readiness probe | `/ready` |
| Service | NodePort |
| Service port | 3000 |
| NodePort | 30080 |
| Service endpoints | 2 pod endpoints |
| Health endpoint | HTTP 200 OK |
| Readiness endpoint | READY |

## Security

The following security controls were implemented:

1. The Docker image runs the application as a non-root user.
2. Kubernetes explicitly uses UID `10001` for the application container.
3. Privilege escalation is disabled.
4. Non-sensitive configuration is separated into a ConfigMap.
5. Sensitive configuration is stored in a Kubernetes Secret.
6. Kubernetes liveness and readiness probes monitor application health.
7. The Sprint 1 application uses a minimal `scratch` runtime image.

## Help

### Common Problems and Solutions

**1. Minikube is not running**

Check the status:

```bash
minikube status
```

Start it with the Docker driver:

```bash
minikube start --driver=docker
```

**2. kubectl cannot connect to Kubernetes**

Check the cluster:

```bash
kubectl get nodes
```

If Minikube is stopped:

```bash
minikube start --driver=docker
```

**3. Application image is not available in Minikube**

Check the local image:

```bash
docker images | grep sprint2-go-app
```

Load it into Minikube:

```bash
minikube image load sprint2-go-app:latest
```

Verify it:

```bash
minikube image ls | grep -i sprint
```

**4. Pods are not running**

Check the pods:

```bash
kubectl get pods
```

Inspect a pod:

```bash
kubectl describe pod <actual-pod-name>
```

For example:

```bash
kubectl describe pod sprint2-go-deployment-5b76b8454c-fl9wv
```

Check logs:

```bash
kubectl logs sprint2-go-deployment-5b76b8454c-fl9wv
```

**5. Health or readiness probe fails**

Inspect the pod:

```bash
kubectl describe pod sprint2-go-deployment-5b76b8454c-fl9wv
```

The probes should show:

```text
Liveness: /health
Readiness: /ready
```

**6. Service endpoints are missing**

Check:

```bash
kubectl get endpoints sprint2-go-service
```

Then check the pods:

```bash
kubectl get pods
```

The Service selector should match the pod label:

```text
app=sprint2-go-app
```

**7. Need the current application URL**

Run:

```bash
minikube service sprint2-go-service --url
```

Do not assume the previous `127.0.0.1` port will remain the same after restarting Minikube.

## Cleanup

The Kubernetes resources can be removed using:

```bash
kubectl delete -f manifests/
```

To stop Minikube:

```bash
minikube stop
```

To completely remove the Minikube cluster:

```bash
minikube delete
```

These cleanup commands were not required for the project submission because the Kubernetes deployment was successfully completed and verified.

### Running the Project

```bash
cd ~/Devops-Intership-2026/sprint-2
minikube start
minikube image load sprint2-go-app:latest
kubectl apply -f manifests/
kubectl get pods
kubectl get services
minikube service sprint2-go-service --url

## Authors

**Anandhrosh**

DevOps Internship 2026

GitHub: https://github.com/Anandhrosh

## Version History

* 1.0
  * Created Kubernetes manifests for the Sprint 1 Go application.
  * Configured a Deployment with two replicas.
  * Added ConfigMap configuration.
  * Added Kubernetes Secret configuration.
  * Added liveness and readiness probes.
  * Added a NodePort Service.
  * Loaded the Docker image into Minikube.
  * Deployed the application to Kubernetes.
  * Verified two running application pods.
  * Verified Service endpoints.
  * Verified `/health` and `/ready` endpoints.
  * Verified pod configuration using `kubectl describe pod`.

## License

This project is intended for educational purposes as part of the DevOps Internship 2026.

## Acknowledgments

* Kubernetes
* Minikube
* Docker
* kubectl
* Git
* GitHub