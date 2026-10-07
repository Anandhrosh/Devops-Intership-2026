# Sprint 3: Kubernetes Package Management with Helm

## Project Title

**Helm Chart for Kubernetes Deployment of a Go Application**

## Description

This project is part of the DevOps Internship 2026 evaluation. The objective of Sprint 3 was to package the Kubernetes application from Sprint 2 into a reusable Helm Chart.

Helm was used to simplify Kubernetes application deployment by converting the Kubernetes manifests into reusable templates and configurable values.

The Helm Chart supports different deployment configurations for development and production environments using separate values files.

The development environment uses one replica with lower resource requirements, while the production environment uses three replicas with higher resource requirements.

The application is deployed to a local Kubernetes cluster using Minikube.

## Project Objectives

* Package the Sprint 2 Kubernetes application into a Helm Chart.
* Create a reusable Helm Chart structure.
* Parameterize the replica count.
* Parameterize the container image repository and tag.
* Parameterize CPU and memory requests and limits.
* Parameterize the Kubernetes Service type and ports.
* Create development environment configuration.
* Create production environment configuration.
* Validate the Helm Chart using `helm lint`.
* Render Kubernetes manifests using `helm template`.
* Deploy the application using `helm install`.
* Verify the Helm release using `helm list` and `helm status`.
* Verify Kubernetes Pods, Deployment, Service, ConfigMap, and Secret.
* Verify application accessibility and health.

## Project Structure

```text
sprint-3/
├── README.md
└── my-app-chart/
    ├── Chart.yaml
    ├── values.yaml
    ├── values-dev.yaml
    ├── values-prod.yaml
    └── templates/
        ├── deployment.yaml
        ├── service.yaml
        ├── configmap.yaml
        └── secret.yaml
```

## Getting Started

### Dependencies

The project was tested using:

* Ubuntu 24.04 WSL
* Docker
* kubectl
* Minikube v1.39.0
* Kubernetes v1.37.0
* Helm v4.3.0
* Git

The Go application image was already created during Sprint 1 and reused during Sprint 2 and Sprint 3.

### Installing

The project was located at:

```text
~/Devops-Intership-2026/sprint-3
```

The Helm Chart was located at:

```text
~/Devops-Intership-2026/sprint-3/my-app-chart
```

The application Docker image was:

```text
sprint2-go-app:latest
```

The image was already available for use with Minikube.

## Executing Program

### Step 1: Verify Helm

The installed Helm version was checked using:

```bash
helm version
```

The installed Helm version was:

```text
v4.3.0
```

### Step 2: Verify Minikube

Minikube was checked using:

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

The Kubernetes node was also verified using:

```bash
kubectl get nodes
```

The node was in `Ready` state and running Kubernetes `v1.37.0`.

### Step 3: Navigate to the Helm Chart

```bash
cd ~/Devops-Intership-2026/sprint-3/my-app-chart
```

### Step 4: Verify the Helm Chart Structure

The Helm Chart contains:

```text
my-app-chart/
├── Chart.yaml
├── values.yaml
├── values-dev.yaml
├── values-prod.yaml
└── templates/
    ├── deployment.yaml
    ├── service.yaml
    ├── configmap.yaml
    └── secret.yaml
```

### Step 5: Validate the Helm Chart

```bash
helm lint .
```

Validation completed successfully:

```text
==> Linting .
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, 0 chart(s) failed
```

The icon message is informational only.

### Step 6: Render the Default Helm Templates

```bash
helm template my-app .
```

The default configuration contains:

```text
Replicas: 2
Image: sprint2-go-app:latest
Service Type: NodePort
Service Port: 3000
NodePort: 30080
```

The Deployment also contains the `/health` liveness probe and `/ready` readiness probe.

### Step 7: Render the Development Configuration

```bash
helm template my-app . -f values-dev.yaml
```

Development configuration:

```text
Replicas: 1

CPU Request: 50m
Memory Request: 32Mi

CPU Limit: 100m
Memory Limit: 64Mi

APP_ENV: development
LOG_LEVEL: debug
```

### Step 8: Render the Production Configuration

```bash
helm template my-app . -f values-prod.yaml
```

Production configuration:

```text
Replicas: 3

CPU Request: 200m
Memory Request: 128Mi

CPU Limit: 500m
Memory Limit: 256Mi

APP_ENV: production
LOG_LEVEL: info
```

The production configuration was validated through Helm template rendering.

## Helm Chart Configuration

### Chart.yaml

The Chart metadata is defined in:

```text
Chart.yaml
```

```yaml
apiVersion: v2
name: my-app
description: A Helm chart for the Sprint 2 Go application
type: application
version: 0.1.0
appVersion: "1.0.0"
```

### values.yaml

The default values file defines the configurable settings for the application.

The main parameters include:

* Replica count
* Container image repository
* Container image tag
* Image pull policy
* Service type
* Service port
* Target port
* NodePort
* CPU requests
* Memory requests
* CPU limits
* Memory limits
* Application environment
* Application name
* Log level

The default replica count is:

```text
2
```

The default Service type is:

```text
NodePort
```

The default NodePort is:

```text
30080
```

### values-dev.yaml

The development values file configures a smaller deployment:

```text
Replicas: 1
CPU Request: 50m
Memory Request: 32Mi
CPU Limit: 100m
Memory Limit: 64Mi
APP_ENV: development
LOG_LEVEL: debug
```

### values-prod.yaml

The production values file configures a larger deployment:

```text
Replicas: 3
CPU Request: 200m
Memory Request: 128Mi
CPU Limit: 500m
Memory Limit: 256Mi
APP_ENV: production
LOG_LEVEL: info
```

## Helm Templates

### Deployment Template

The Deployment template is located at:

```text
templates/deployment.yaml
```

The Deployment uses Helm values for:

* Replica count
* Container image
* Image pull policy
* CPU requests
* Memory requests
* CPU limits
* Memory limits

The application also uses Kubernetes health probes.

### Liveness Probe

The liveness probe checks:

```text
/health
```

The probe is configured on port:

```text
3000
```

The liveness probe allows Kubernetes to determine whether the application is running correctly.

### Readiness Probe

The readiness probe checks:

```text
/ready
```

The probe is configured on port:

```text
3000
```

The readiness probe allows Kubernetes to determine whether the application is ready to receive traffic.

### Container Security Configuration

The Deployment uses:

```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 10001
  allowPrivilegeEscalation: false
```

This ensures that the application container runs as a non-root user and prevents privilege escalation.

### Service Template

The Service template is located at:

```text
templates/service.yaml
```

The default Service configuration is:

```text
Service Type: NodePort
Service Port: 3000
Target Port: 3000
NodePort: 30080
```

The template conditionally configures the NodePort when the Service type is `NodePort`.

### ConfigMap Template

The ConfigMap template is located at:

```text
templates/configmap.yaml
```

The ConfigMap contains:

```text
APP_ENV
APP_NAME
LOG_LEVEL
```

These values are controlled by the Helm configuration files.

### Secret Template

The Secret template is located at:

```text
templates/secret.yaml
```

The Secret provides the demonstration application secret used by the Deployment.

The Secret in this project is only a demonstration value and is not a real production credential.

Real production secrets should not be committed directly to Git.

## Installing the Helm Release

### Step 9: Install the Development Environment

```bash
helm install my-app-dev . -f values-dev.yaml
```

The installation completed successfully:

```text
NAME: my-app-dev
NAMESPACE: default
STATUS: deployed
REVISION: 1
DESCRIPTION: Install complete
```

### Step 10: Verify Helm Release

Check Helm releases:

```bash
helm list
```

The release was shown as:

```text
NAME          NAMESPACE   REVISION   STATUS
my-app-dev    default     1          deployed
```

Check release details:

```bash
helm status my-app-dev
```

The release status was:

```text
STATUS: deployed
```

## Kubernetes Verification

### Step 11: Verify the Pods

```bash
kubectl get pods
```

The deployed Pod was:

```text
sprint2-go-deployment-5698f7c7bb-g5vzw   1/1   Running   0
```

The Pod was ready and had zero restarts.

### Step 12: Verify the Deployment

```bash
kubectl get deployment
```

The Deployment state was:

```text
NAME                    READY   UP-TO-DATE   AVAILABLE
sprint2-go-deployment   1/1     1            1
```

The development environment successfully deployed one replica as configured in `values-dev.yaml`.

### Step 13: Verify the Service

```bash
kubectl get service
```

The application Service was:

```text
sprint2-go-service   NodePort   10.98.218.176   <none>   3000:30080/TCP
```

The Service configuration was:

```text
Service type: NodePort
Service port: 3000
Target port: 3000
NodePort: 30080
```

### Step 14: Verify the ConfigMap

```bash
kubectl get configmap sprint2-config -o yaml
```

The development values were successfully applied:

```text
APP_ENV: development
APP_NAME: sprint2-go-app
LOG_LEVEL: debug
```

The ConfigMap also contained Helm ownership metadata, confirming that it was created and managed by the Helm release.

### Step 15: Verify the Deployment Configuration

```bash
kubectl get deployment sprint2-go-deployment -o yaml
```

The Deployment confirmed:

```text
Replicas: 1
Image: sprint2-go-app:latest
CPU Request: 50m
Memory Request: 32Mi
CPU Limit: 100m
Memory Limit: 64Mi
Liveness Probe: /health
Readiness Probe: /ready
Run as non-root: true
Run as user: 10001
Privilege escalation: false
```

### Step 16: Get the Application URL

The application was exposed through Minikube using:

```bash
minikube service sprint2-go-service --url
```

During the verification session, Minikube returned temporary port like this:

```text
http://127.0.0.1:45843
```

The URL is generated by Minikube for the active session and may change after restarting Minikube or the service.

Because the Docker driver is being used on Linux/WSL, the terminal running the `minikube service --url` command must remain open.

## Application Testing

### Step 17: Test the Application Home Page

The application was tested using:

```bash
curl -v http://127.0.0.1:45843/
```

The application returned:

```text
HTTP/1.1 200 OK
```

This confirmed that the Go application was successfully reachable through the Kubernetes Service.

### Step 18: Test the Health Endpoint

The health endpoint was tested using:

```bash
curl -v http://127.0.0.1:45843/health
```

The application returned:

```text
HTTP/1.1 200 OK
```

The response displayed the application health status as healthy.

The same `/health` endpoint is used by the Kubernetes liveness probe.

## Development and Production Configuration

### Development Environment

The development environment uses:

```text
Replicas: 1
CPU Request: 50m
Memory Request: 32Mi
CPU Limit: 100m
Memory Limit: 64Mi
APP_ENV: development
LOG_LEVEL: debug
```

Render the development configuration:

```bash
helm template my-app . -f values-dev.yaml
```

Install the development configuration:

```bash
helm install my-app-dev . -f values-dev.yaml
```

### Production Environment

The production environment uses:

```text
Replicas: 3
CPU Request: 200m
Memory Request: 128Mi
CPU Limit: 500m
Memory Limit: 256Mi
APP_ENV: production
LOG_LEVEL: info
```

Render the production configuration:

```bash
helm template my-app . -f values-prod.yaml
```

The production configuration was validated through Helm template rendering.

## Verification Summary

The following checks were completed successfully:

| Check | Result |
|---|---|
| Helm installation | v4.3.0 |
| Helm lint | Passed |
| Default template rendering | Passed |
| Development template rendering | Passed |
| Production template rendering | Passed |
| Development Helm installation | Passed |
| Helm release status | Deployed |
| Development replicas | 1 |
| Application Pod | 1/1 Running |
| Pod restarts | 0 |
| ConfigMap | Created and referenced |
| Secret | Created and referenced |
| Liveness probe | `/health` |
| Readiness probe | `/ready` |
| Service | NodePort |
| Service port | 3000 |
| NodePort | 30080 |
| Application home page | HTTP 200 OK |
| Health endpoint | HTTP 200 OK |

## Security

The following security controls were maintained from the Sprint 2 application:

1. The Docker image runs the application as a non-root user.
2. Kubernetes explicitly uses UID `10001` for the application container.
3. Privilege escalation is disabled.
4. Non-sensitive configuration is stored in a ConfigMap.
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

Start Minikube:

```bash
minikube start --driver=docker
```

**2. Helm release is not deployed**

Check Helm releases:

```bash
helm list
```

Check the release:

```bash
helm status my-app-dev
```

**3. Helm Chart validation fails**

Run:

```bash
helm lint .
```

Render the templates:

```bash
helm template my-app .
```

These commands help identify errors in the Helm Chart before deployment.

**4. Pod is not running**

Check:

```bash
kubectl get pods
```

Inspect the Pod:

```bash
kubectl describe pod <pod-name>
```

Check application logs:

```bash
kubectl logs <pod-name>
```

**5. Service cannot be accessed**

Check the Service:

```bash
kubectl get service
```

Then obtain the current Minikube URL:

```bash
minikube service sprint2-go-service --url
```

Keep the terminal running while using the generated URL.

**6. Need to check the current Helm configuration**

Development:

```bash
helm template my-app . -f values-dev.yaml
```

Production:

```bash
helm template my-app . -f values-prod.yaml
```

## Helm Commands Reference

### Validate Chart

```bash
helm lint .
```

### Render Templates

```bash
helm template my-app .
```

### Render Development Configuration

```bash
helm template my-app . -f values-dev.yaml
```

### Render Production Configuration

```bash
helm template my-app . -f values-prod.yaml
```

### Install

```bash
helm install my-app-dev . -f values-dev.yaml
```

### List Releases

```bash
helm list
```

### Check Release Status

```bash
helm status my-app-dev
```

### Upgrade Release

```bash
helm upgrade my-app-dev . -f values-dev.yaml
```

### Uninstall Release

```bash
helm uninstall my-app-dev
```

## Cleanup

To remove the Helm release:

```bash
helm uninstall my-app-dev
```

To stop Minikube:

```bash
minikube stop
```

To completely remove the Minikube cluster:

```bash
minikube delete
```

The cleanup commands are not required during normal testing because the Helm release can remain deployed for verification.

## Run the Project

### 1. Go to the Helm chart directory

```bash
cd ~/Devops-Intership-2026/sprint-3/my-app-chart
```
2. Start Minikube
```bash
minikube start
```
3. Check Minikube status
```bash
minikube status
```
4. Install the Helm chart
```bash
helm install my-app-dev . -f values-dev.yaml
```
5. Check Helm release
```bash
helm list
```
6. Check Kubernetes resources
```bash
kubectl get pods
kubectl get deployment
kubectl get service
```
7. Get the application URL
```bash
minikube service sprint2-go-service --url
```
Example:

http://127.0.0.1:45843

The local port may be different each time. This is normal because Minikube creates a temporary local tunnel.

8. Test the application

Open another terminal and use the URL returned above:

curl http://127.0.0.1:45843/

Test the health endpoint:

curl http://127.0.0.1:45843/health

9. Check Helm status
```bash
helm status my-app-dev
```
10. Stop and remove the project
```bash
helm uninstall my-app-dev
minikube stop
```

## Authors

**Anandhrosh**

DevOps Internship 2026

GitHub: https://github.com/Anandhrosh

## Version History

* 1.0
  * Created a reusable Helm Chart for the Sprint 2 Go application.
  * Added `Chart.yaml`.
  * Added default `values.yaml`.
  * Added development environment values.
  * Added production environment values.
  * Created Helm Deployment template.
  * Created Helm Service template.
  * Created Helm ConfigMap template.
  * Created Helm Secret template.
  * Parameterized replica count.
  * Parameterized container image repository and tag.
  * Parameterized CPU and memory resources.
  * Parameterized Service type and ports.
  * Validated the Chart using `helm lint`.
  * Rendered default, development, and production templates.
  * Installed the development Helm release.
  * Verified the Helm release using `helm list` and `helm status`.
  * Verified the Kubernetes Pod and Deployment.
  * Verified the NodePort Service.
  * Verified the application home page with HTTP 200.
  * Verified the `/health` endpoint with HTTP 200.

## License

This project is intended for educational purposes as part of the DevOps Internship 2026.

## Acknowledgments

* Kubernetes
* Helm
* Minikube
* Docker
* kubectl
* Git
* GitHub
