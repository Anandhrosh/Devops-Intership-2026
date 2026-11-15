# Sprint 10: Three-Tier Kubernetes Capstone Project

## Project Title

**Three-Tier Microservices Application Deployment using Kubernetes**

## Description

This project demonstrates the deployment of a three-tier web application on a local Kubernetes cluster using Minikube.

The application consists of a React frontend, a Go backend REST API, and a PostgreSQL database. Users can upload images, view uploaded images in a gallery, and delete images. The backend stores image data in PostgreSQL using the BYTEA data type.

The project implements Kubernetes Deployments, Services, StatefulSets, Persistent Volume Claims, Secrets, Horizontal Pod Autoscaling, NGINX Ingress, NetworkPolicy, and Kustomize.

The backend runs multiple replicas with liveness and readiness probes. The frontend is exposed through NGINX Ingress using the host `capstone.local`. PostgreSQL uses persistent storage to retain uploaded images.

An ArgoCD Application manifest is also prepared to demonstrate GitOps configuration. ArgoCD installation and synchronization remain pending because of the resource limitations of the local environment.

## Project Objectives

- Develop a React frontend application.
- Develop a Go backend REST API.
- Integrate the backend with PostgreSQL.
- Implement image upload, gallery display, and image deletion.
- Containerize frontend and backend applications using Docker.
- Deploy the frontend and backend using Kubernetes Deployments.
- Deploy PostgreSQL using a StatefulSet.
- Configure Persistent Volume Claims for database storage.
- Store database credentials using Kubernetes Secrets.
- Configure backend liveness and readiness probes.
- Configure Horizontal Pod Autoscaling for the backend.
- Expose the frontend through NGINX Ingress.
- Configure a NetworkPolicy to restrict database ingress.
- Package Kubernetes resources using Kustomize.
- Prepare an ArgoCD Application manifest for GitOps.
- Verify application functionality and Kubernetes resources.
- Document implementation and verification using screenshots.

## Getting Started

### Dependencies

Ensure the following tools are available in your Ubuntu/WSL environment:

- Docker Desktop with WSL integration
- Ubuntu on WSL
- Minikube
- Kubernetes `kubectl`
- Git
- Bash
- Node.js and npm (for frontend development)
- Go (for backend development)

### Installing

1. Open the Ubuntu terminal and move to the repository:

   ```bash
   cd ~/Devops-Intership-2026
   ```

2. Start Minikube:

   ```bash
   minikube start
   ```

3. Confirm that the Kubernetes node is ready:

   ```bash
   kubectl get nodes
   ```

   Expected result: the Minikube node should show `Ready`.

4. Check the Minikube status:

   ```bash
   minikube status
   ```

5. Move to the Sprint 10 project directory:

   ```bash
   cd ~/Devops-Intership-2026/sprint-10
   ```

6. Verify the project structure:

   ```bash
   ls -R
   ```

7. Verify that the Kubernetes resources are available:

   ```bash
   find kubernetes -maxdepth 2 -type f | sort
   ```

8. Verify the frontend and backend Dockerfiles:

   ```bash
   ls -l frontend/Dockerfile backend/Dockerfile
   ```

## Executing Program

### 1. Create the Sprint 10 directory structure

The project is organized into frontend, backend, Kubernetes, and ArgoCD directories.

```bash
cd ~/Devops-Intership-2026/sprint-10
```

The directory structure is:

```text
sprint-10/
├── README.md
├── frontend/
│   ├── Dockerfile
│   └── application/
├── backend/
│   ├── Dockerfile
│   └── application/
├── kubernetes/
│   ├── database/
│   │   ├── pvc.yaml
│   │   ├── service.yaml
│   │   ├── statefulset.yaml
│   │   └── kustomization.yaml
│   ├── backend/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   ├── hpa.yaml
│   │   └── kustomization.yaml
│   ├── frontend/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── kustomization.yaml
│   ├── ingress.yaml
│   ├── network-policy.yaml
│   └── kustomization.yaml
└── argocd/
    └── application.yaml
```

### 2. Build the frontend Docker image

The frontend is developed using React and Vite.

The Dockerfile uses a multi-stage build. The first stage builds the React application, and the second stage serves the generated files using an unprivileged NGINX image.

Move to the frontend directory:

```bash
cd ~/Devops-Intership-2026/sprint-10/frontend
```

Build the image:

```bash
docker build -t capstone-frontend:1.0 .
```

Verify the image:

```bash
docker images | grep capstone-frontend
```

### 3. Build the backend Docker image

The backend is developed using Go and provides REST API endpoints for image operations and application status.

The Dockerfile uses a multi-stage build and a Distroless runtime image.

Move to the backend directory:

```bash
cd ~/Devops-Intership-2026/sprint-10/backend
```

Build the image:

```bash
docker build -t capstone-backend:1.0 .
```

Verify the image:

```bash
docker images | grep capstone-backend
```

### 4. Load Docker images into Minikube

Load the frontend image:

```bash
minikube image load capstone-frontend:1.0
```

Load the backend image:

```bash
minikube image load capstone-backend:1.0
```

Verify that the images are available:

```bash
minikube image ls | grep capstone
```

### 5. Configure the PostgreSQL Secret

The database credentials are stored in a Kubernetes Secret named:

`capstone-postgres-secret`

The Secret contains the following keys:

- `POSTGRES_DB`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`

The database name is `capstonedb`, and the database user is `capstoneuser`.

The plaintext Secret YAML is intentionally excluded from Git to protect database credentials.

For a fresh deployment, create or provision the Secret securely before applying the database and backend resources. Use credentials consistent with the PostgreSQL database initialization.

**Important:** The existing cluster already has this Secret. Do not replace it casually because the database uses persistent storage and its existing credentials must remain synchronized.

### 6. Deploy the PostgreSQL StatefulSet

PostgreSQL is deployed using a StatefulSet to provide stable workload identity and persistent storage.

The database configuration includes:

- StatefulSet: `capstone-postgres`
- Replicas: 1
- Image: `postgres:17-alpine`
- Service: `capstone-postgres-service`
- PVC: `capstone-postgres-pvc`
- Storage: 5Gi
- Database: `capstonedb`

Apply the database resources:

```bash
kubectl apply -k kubernetes/database/
```

Verify the StatefulSet:

```bash
kubectl get statefulsets
```

Verify the PostgreSQL Pod:

```bash
kubectl get pods -l app=capstone-postgres
```

Verify persistent storage:

```bash
kubectl get pvc capstone-postgres-pvc
```

Expected result: the PostgreSQL Pod should be `Running`, and the PVC should show `Bound`.

### 7. Deploy the Go backend

The backend Deployment runs two replicas initially.

The configuration includes:

- Deployment: `capstone-backend`
- Replicas: 2
- Service: `capstone-backend-service`
- Container port: 3000
- Liveness endpoint: `/health`
- Readiness endpoint: `/ready`
- Database connection through PostgreSQL Service
- Database credentials loaded from Kubernetes Secret

Apply the backend resources:

```bash
kubectl apply -k kubernetes/backend/
```

Verify the Deployment:

```bash
kubectl get deployments capstone-backend
```

Verify the backend Pods:

```bash
kubectl get pods -l app=capstone-backend
```

Verify the backend Service:

```bash
kubectl get service capstone-backend-service
```

Expected result: two backend replicas should become ready.

### 8. Configure Horizontal Pod Autoscaling

The backend uses a Horizontal Pod Autoscaler to adjust the number of replicas based on CPU utilization.

The HPA configuration includes:

- HPA: `capstone-backend-hpa`
- Minimum replicas: 2
- Maximum replicas: 5
- CPU utilization target: 50%

Apply the backend HPA through the backend Kustomize configuration.

Verify the HPA:

```bash
kubectl get hpa capstone-backend-hpa
```

Check detailed autoscaling information:

```bash
kubectl describe hpa capstone-backend-hpa
```

Expected result: the HPA should display the current CPU utilization and the configured target.

### 9. Deploy the React frontend

The frontend is deployed using a Kubernetes Deployment.

The configuration includes:

- Deployment: `capstone-frontend`
- Replicas: 2
- Service: `capstone-frontend-service`
- Container port: 8080
- Liveness probe: `/`
- Readiness probe: `/`

Apply the frontend resources:

```bash
kubectl apply -k kubernetes/frontend/
```

Verify the frontend Deployment:

```bash
kubectl get deployments capstone-frontend
```

Verify the frontend Pods:

```bash
kubectl get pods -l app=capstone-frontend
```

Verify the frontend Service:

```bash
kubectl get service capstone-frontend-service
```

Expected result: both frontend replicas should become ready.

### 10. Configure NGINX Ingress

NGINX Ingress provides external access to the frontend and backend services.

The Ingress resource is named:

`capstone-ingress`

The configured host is:

`capstone.local`

Routing configuration:

| Path | Destination |
|---|---|
| `/` | `capstone-frontend-service:8080` |
| `/api` | `capstone-backend-service:3000` |

Enable the Minikube Ingress addon if it is not already enabled:

```bash
minikube addons enable ingress
```

Verify the Ingress controller:

```bash
kubectl get pods -n ingress-nginx
```

Apply the Ingress:

```bash
kubectl apply -f kubernetes/ingress.yaml
```

Verify the Ingress resource:

```bash
kubectl get ingress capstone-ingress
```

### 11. Access the application

The application is configured to use the hostname:

`capstone.local`

For the existing WSL and Minikube environment, port forwarding is used to access the Ingress controller.

Run:

```bash
kubectl port-forward --address 0.0.0.0 -n ingress-nginx service/ingress-nginx-controller 8088:80
```

Keep this terminal running.

On Windows, add the following entry to the hosts file using the current WSL IP address:

```text
<WSL-IP-ADDRESS> capstone.local
```

The WSL IP can be checked using:

```bash
hostname -I
```

Open the browser:

```text
http://capstone.local:8088
```

The React application should load.

### 12. Configure NetworkPolicy

The project includes a NetworkPolicy named:

`capstone-database-policy`

The policy selects PostgreSQL Pods using:

```yaml
app: capstone-postgres
```

It allows incoming PostgreSQL traffic on TCP port `5432` from Pods labeled:

```yaml
tier: backend
```

Apply the policy:

```bash
kubectl apply -f kubernetes/network-policy.yaml
```

Verify the policy:

```bash
kubectl get networkpolicy
```

**Important limitation:** The current Minikube cluster uses the `kindnet` CNI, which does not enforce Kubernetes NetworkPolicies. The policy manifest is created and applied, but actual network isolation has not been demonstrated. A NetworkPolicy-compatible CNI is required for enforcement.

### 13. Configure Kustomize

Kustomize is used to organize Kubernetes resources into reusable configurations.

The project contains separate Kustomize configurations for:

- Database
- Backend
- Frontend
- Root application resources

Preview the complete configuration:

```bash
kubectl kustomize kubernetes/
```

Check differences against the running cluster:

```bash
kubectl diff -k kubernetes/
```

Apply the complete configuration:

```bash
kubectl apply -k kubernetes/
```

The database Secret is intentionally excluded from Kustomize and must be provisioned separately.

### 14. Prepare ArgoCD GitOps configuration

An ArgoCD Application manifest is prepared at:

```text
argocd/application.yaml
```

The manifest references:

- Repository: `https://github.com/Anandhrosh/Devops-Intership-2026.git`
- Branch: `main`
- Application path: `sprint-10/kubernetes`
- Destination namespace: `default`

The manifest is prepared for GitOps synchronization.

**Current status:** ArgoCD has not been installed in the Minikube cluster because of local memory constraints. Application synchronization has not yet been performed.

The database Secret also requires separate provisioning because it is not stored in Git.

### 15. Final verification commands

Run the following commands from the Sprint 10 directory:

```bash
kubectl get deployments
kubectl get pods -o wide
kubectl get services
kubectl get statefulsets
kubectl get pvc
kubectl get hpa
kubectl get ingress
kubectl get networkpolicy
```

Check backend application status through Ingress:

```bash
minikube ssh -- "curl -s --max-time 10 -H 'Host: capstone.local' http://127.0.0.1:31352/api/status"
```

Expected application response:

```json
{
  "application": "DevOps Capstone",
  "backend": "Running",
  "database": "Connected"
}
```

These commands verify the deployed application resources and backend database connectivity.

## Technical Explanation

### What is a Three-Tier Application?

A three-tier application separates the system into three main layers:

1. **Presentation Layer:** The frontend that users interact with.
2. **Application Layer:** The backend that processes requests and business logic.
3. **Data Layer:** The database that stores application data.

In this project, React provides the presentation layer, Go provides the application layer, and PostgreSQL provides the data layer.

### What is React?

React is a JavaScript library used to build user interfaces.

In this project, React provides the image upload interface, image gallery, and image deletion functionality.

### What is Go?

Go is a programming language used to build backend services and APIs.

The Go backend receives requests from the frontend, communicates with PostgreSQL, and provides image upload, retrieval, deletion, health, and status endpoints.

### What is PostgreSQL?

PostgreSQL is a relational database management system.

The project uses PostgreSQL to store image metadata and image content. Image data is stored using the BYTEA data type.

### What is a Kubernetes Deployment?

A Deployment manages stateless application Pods and maintains the desired number of replicas.

In this project, the frontend and backend use Deployments to maintain two replicas initially.

### What is a StatefulSet?

A StatefulSet manages stateful applications that require stable identity and persistent storage.

In this project, PostgreSQL uses a StatefulSet named `capstone-postgres`.

### What is a Persistent Volume Claim?

A Persistent Volume Claim requests persistent storage from Kubernetes.

The project uses a 5Gi PVC named `capstone-postgres-pvc` to retain PostgreSQL data independently of the PostgreSQL Pod lifecycle.

### What is a Kubernetes Secret?

A Secret stores sensitive configuration data such as database usernames and passwords.

The project uses `capstone-postgres-secret` to provide PostgreSQL credentials to the database and backend containers.

### What is a Kubernetes Service?

A Service provides a stable network endpoint for accessing a group of Pods.

The project uses ClusterIP Services for the frontend, backend, and PostgreSQL components.

### What is Horizontal Pod Autoscaling?

Horizontal Pod Autoscaling automatically adjusts the number of Pod replicas based on configured resource metrics.

The backend HPA is configured with a minimum of two replicas, a maximum of five replicas, and a CPU utilization target of 50%.

### What are Liveness and Readiness Probes?

- **Liveness probe:** Checks whether a container is healthy and should continue running.
- **Readiness probe:** Checks whether a container is ready to receive traffic.

The backend uses `/health` for liveness and `/ready` for readiness.

### What is NGINX Ingress?

NGINX Ingress manages external HTTP traffic entering the Kubernetes cluster.

In this project, it routes frontend requests to the React Service and API requests to the Go backend Service using the hostname `capstone.local`.

### What is a NetworkPolicy?

A NetworkPolicy defines rules that control network communication between Kubernetes Pods.

The project includes a policy intended to allow database ingress only from backend Pods on TCP port 5432.

Actual enforcement requires a compatible CNI plugin.

### What is Kustomize?

Kustomize is a Kubernetes configuration management tool that organizes and customizes YAML manifests.

The project uses separate Kustomize configurations for the frontend, backend, and database, along with a root configuration.

### What is GitOps?

GitOps is an approach in which Git repositories store the desired application configuration, and a controller continuously reconciles the cluster with that configuration.

The project includes an ArgoCD Application manifest, but installation and synchronization remain pending.

## Security and Resource Management

- Frontend and backend containers run as non-root users.
- Privilege escalation is disabled.
- Linux capabilities are dropped.
- Read-only root filesystems are configured where applicable.
- Database credentials are provided through Kubernetes Secrets.
- Plaintext database credentials are excluded from Git.
- Frontend and backend resource requests and limits are configured.
- Backend health and readiness probes are enabled.
- PostgreSQL data is stored using persistent storage.
- A NetworkPolicy manifest defines the intended database ingress restriction.

## Help

### Minikube is not running

Check the cluster:

```bash
minikube status
kubectl get nodes
```

Start Minikube if necessary:

```bash
minikube start
```

### Docker image is not available in Minikube

Check the images:

```bash
minikube image ls | grep capstone
```

Load the images if necessary:

```bash
minikube image load capstone-frontend:1.0
minikube image load capstone-backend:1.0
```

### A Pod is not running

Check the Pods:

```bash
kubectl get pods
```

Inspect a Pod:

```bash
kubectl describe pod <pod-name>
```

Check container logs:

```bash
kubectl logs <pod-name>
```

### Backend cannot connect to PostgreSQL

Check the database Pod:

```bash
kubectl get pods -l app=capstone-postgres
```

Check the database Service:

```bash
kubectl get service capstone-postgres-service
```

Check the Secret exists without displaying its contents:

```bash
kubectl get secret capstone-postgres-secret
```

Inspect backend logs:

```bash
kubectl logs deployment/capstone-backend
```

### Persistent Volume Claim is not bound

Check storage:

```bash
kubectl get pvc
kubectl get pv
```

Inspect the claim:

```bash
kubectl describe pvc capstone-postgres-pvc
```

### HPA does not show CPU metrics

Check metrics-server:

```bash
kubectl get pods -n kube-system | grep metrics-server
```

Check resource metrics:

```bash
kubectl top pods
```

Check the HPA:

```bash
kubectl get hpa
```

### Ingress is not accessible

Check the Ingress controller:

```bash
kubectl get pods -n ingress-nginx
```

Check the Ingress Service:

```bash
kubectl get service -n ingress-nginx
```

Check the Ingress resource:

```bash
kubectl get ingress capstone-ingress
```

Check that port forwarding is running and that the Windows hosts file contains the correct current WSL IP address.

### Application is not accessible from Windows

Verify the port-forward process:

```bash
ss -lntp | grep 8088
```

Check the application from inside Minikube:

```bash
minikube ssh -- "curl -I --max-time 5 -H 'Host: capstone.local' http://127.0.0.1:31352/"
```

### NetworkPolicy is not enforcing restrictions

The current cluster uses `kindnet`, which does not enforce NetworkPolicies.

A compatible CNI plugin is required before network isolation can be verified.

### ArgoCD is not available

ArgoCD is not currently installed in the cluster.

The Application manifest is available in `argocd/application.yaml`. Install ArgoCD and configure the repository connection before attempting synchronization.

## Run the Project

### 1. Go to the Sprint 10 directory

```bash
cd ~/Devops-Intership-2026/sprint-10
```

### 2. Start Minikube

```bash
minikube start
```

### 3. Verify the cluster

```bash
kubectl get nodes
```

### 4. Build the frontend and backend images

```bash
cd frontend
docker build -t capstone-frontend:1.0 .
cd ../backend
docker build -t capstone-backend:1.0 .
```

### 5. Load the images into Minikube

```bash
minikube image load capstone-frontend:1.0
minikube image load capstone-backend:1.0
```

### 6. Verify the database Secret

```bash
kubectl get secret capstone-postgres-secret
```

For a fresh cluster, provision the Secret securely before deploying the database and backend.

### 7. Deploy Kubernetes resources

From the Sprint 10 directory:

```bash
cd ~/Devops-Intership-2026/sprint-10
kubectl apply -k kubernetes/
```

### 8. Verify application resources

```bash
kubectl get deployments,pods,services
kubectl get statefulsets,pvc,hpa
kubectl get ingress,networkpolicy
```

### 9. Access the application

Start the Ingress port-forward:

```bash
kubectl port-forward --address 0.0.0.0 -n ingress-nginx service/ingress-nginx-controller 8088:80
```

Open:

```text
http://capstone.local:8088
```

## Proof of Work

The following screenshots document the implementation and verification of the Sprint 10 capstone.

Save the screenshots in a `Screenshots/` directory inside `sprint-10/` using the suggested filenames below.

### 1. Kubernetes Cluster and Deployments

This screenshot should show the Minikube node and frontend/backend Deployments running successfully.

Suggested filename:

`Screenshots/01-kubernetes-deployments.png`

### 2. PostgreSQL StatefulSet and Persistent Storage

This screenshot should show the PostgreSQL StatefulSet, running Pod, and Bound 5Gi PVC.

Suggested filename:

`Screenshots/02-postgresql-persistent-storage.png`

### 3. Backend HPA

This screenshot should show the backend HPA with its minimum and maximum replica configuration and CPU target.

Suggested filename:

`Screenshots/03-backend-hpa.png`

### 4. NGINX Ingress

This screenshot should show the `capstone-ingress` resource and its `capstone.local` hostname.

Suggested filename:

`Screenshots/04-nginx-ingress.png`

### 5. Application Running in Browser

This screenshot should show the React application opened at `http://capstone.local:8088`.

Suggested filename:

`Screenshots/05-application-browser.png`

### 6. Image Upload and Gallery

This screenshot should demonstrate successful image upload and display in the gallery.

Suggested filename:

`Screenshots/06-image-upload-gallery.png`

### 7. Image Deletion

This screenshot should demonstrate the image deletion functionality.

Suggested filename:

`Screenshots/07-image-deletion.png`

### 8. Backend and Database Status

This screenshot should show the backend status response confirming that the backend is running and the database is connected.

Suggested filename:

`Screenshots/08-backend-database-status.png`

### 9. Kustomize Configuration

This screenshot should show the complete Kubernetes resource configuration rendered using Kustomize.

Suggested filename:

`Screenshots/09-kustomize-resources.png`

### 10. NetworkPolicy Manifest

This screenshot should show the database NetworkPolicy configuration.

Suggested filename:

`Screenshots/10-network-policy.png`

**Note:** The NetworkPolicy manifest is present, but enforcement is unavailable with the current `kindnet` CNI.

### 11. ArgoCD Application Manifest

This screenshot should show the prepared ArgoCD Application YAML referencing the GitHub repository and Sprint 10 Kubernetes path.

Suggested filename:

`Screenshots/11-argocd-application.png`

**Note:** ArgoCD installation and synchronization are pending.

## Authors

**Anandhrosh S**

DevOps Internship 2026

GitHub: https://github.com/Anandhrosh

Repository: https://github.com/Anandhrosh/Devops-Intership-2026

## Version History

- **0.1.0** — Initial Sprint 10 implementation with React frontend, Go backend, PostgreSQL StatefulSet, persistent storage, Kubernetes Deployments, HPA, NGINX Ingress, Kustomize configuration, and ArgoCD Application manifest preparation.

## License

This project is intended for educational purposes as part of the DevOps Internship 2026.

## Acknowledgments

- Kubernetes documentation
- Minikube documentation
- Docker documentation
- React documentation
- Go documentation
- PostgreSQL documentation
- Kustomize documentation
- ArgoCD documentation