# Sprint 4: Ingress Routing & SSL/TLS Certificate Management

## Project Title

**Ingress Routing & SSL/TLS Certificate Management**

## Description

This project demonstrates how to expose a Kubernetes application through
an NGINX Ingress Controller and secure HTTP traffic with HTTPS using a
TLS certificate.

The existing Sprint 2 Go application is exposed through the Kubernetes
Service `sprint2-go-service`. An Ingress resource routes requests for
the hostname `app.local` to that Service. A Kubernetes TLS Secret stores
the certificate and private key used by the Ingress Controller.

For this local lab, a self-signed certificate is used. Browsers and
`curl` do not automatically trust self-signed certificates, so the test
commands use `-k` to skip certificate trust verification. This is
suitable for a learning environment, not a production deployment.

## Project Objectives

-   Enable the NGINX Ingress Controller in Minikube.
-   Create a self-signed TLS certificate for `app.local`.
-   Store the certificate and private key in a Kubernetes TLS Secret.
-   Configure an Ingress resource to route requests to
    `sprint2-go-service`.
-   Enable HTTPS and redirect HTTP traffic to HTTPS.
-   Verify that the application can be accessed through the configured
    hostname.

## Getting Started

### Dependencies

Ensure the following tools are available in your Ubuntu/WSL environment:

-   Docker Desktop with WSL integration
-   Ubuntu on WSL
-   Minikube
-   kubectl
-   OpenSSL
-   curl

### Installing

1.  Open the Ubuntu terminal and move to the repository:

    ``` bash
    cd ~/Devops-Intership-2026
    ```

2.  Start Minikube:

    ``` bash
    minikube start
    ```

3.  Confirm that the Kubernetes node is ready:

    ``` bash
    kubectl get nodes
    ```

    Expected result: the Minikube node should show `Ready`.

4.  Confirm that the Sprint 2 application Service exists:

    ``` bash
    kubectl get deployment sprint2-go-deployment
    kubectl get service sprint2-go-service
    ```

    The Ingress configuration in this project forwards requests to
    `sprint2-go-service` on Service port `3000`. Deploy the Sprint 2
    application first if these resources are missing.

## Executing Program

### 1. Create the Sprint 4 directory

``` bash
mkdir -p sprint-4/ingress
cd sprint-4/ingress
```

### 2. Enable the NGINX Ingress Controller

``` bash
minikube addons enable ingress
```

Verify that the controller is running:

``` bash
kubectl get pods -n ingress-nginx
kubectl get service -n ingress-nginx ingress-nginx-controller
```

The controller should reach `Running` status. Minikube may expose the
controller through NodePorts for HTTP and HTTPS.

### 3. Generate a self-signed TLS certificate

Run this command from `sprint-4/ingress`:

``` bash
openssl req -x509 -nodes -days 365 \
  -newkey rsa:2048 \
  -keyout tls.key \
  -out tls.crt \
  -subj "/CN=app.local" \
  -addext "subjectAltName=DNS:app.local"
```

This creates:

-   `tls.crt` --- the certificate.
-   `tls.key` --- the private key. Keep this file secret and do not
    commit it to Git.

The certificate is self-signed and intended for this local exercise.

### 4. Create the Kubernetes TLS Secret

``` bash
kubectl create secret tls app-local-tls --cert=tls.crt --key=tls.key
```

Verify the Secret:

``` bash
kubectl get secret app-local-tls
kubectl describe secret app-local-tls
```

The Secret type should be `kubernetes.io/tls` and it should contain the
`tls.crt` and `tls.key` data entries. Do not print or share the private
key contents.

### 5. Configure the Ingress resource

The file `sprint-4/ingress/ingress.yaml` contains:

``` yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-local-ingress
  namespace: default
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - app.local
      secretName: app-local-tls
  rules:
    - host: app.local
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: sprint2-go-service
                port:
                  number: 3000
```

Apply the manifest from the repository root:

``` bash
cd ~/Devops-Intership-2026
kubectl apply -f sprint-4/ingress/ingress.yaml
```

Verify the Ingress:

``` bash
kubectl get ingress app-local-ingress
kubectl describe ingress app-local-ingress
```

The Ingress should list `app.local` as its host and use the
`app-local-tls` Secret for TLS.

### 6. Test HTTPS from Ubuntu/WSL

If direct access to the Minikube IP or NodePort is unavailable in your
WSL setup, use port-forwarding as a local testing workaround.

In one Ubuntu terminal, run:

``` bash
kubectl port-forward --address 0.0.0.0 \
  -n ingress-nginx service/ingress-nginx-controller 8443:443
```

Keep this terminal open. In a second Ubuntu terminal, run:

``` bash
curl -k -v --resolve app.local:8443:127.0.0.1 https://app.local:8443/
```

A successful response should return HTTP `200` and the Go application's
HTML. The `-k` option skips certificate trust verification because the
certificate is self-signed.

### 7. Optional: Test from Windows

The following steps are specific to the Windows + WSL local lab. They
are not required when testing from Ubuntu.

1.  Keep the WSL port-forward command from Step 6 running.

2.  Open **PowerShell as Administrator** and add a Windows port proxy,
    replacing the WSL IP below if it has changed:

    ``` powershell
    netsh interface portproxy add v4tov4 listenaddress=127.0.0.1 listenport=443 connectaddress=172.28.58.73 connectport=8443
    ```

3.  Open the Windows hosts file as Administrator:

    `C:\Windows\System32\drivers\etc\hosts`

    Add this line:

    ``` text
    127.0.0.1 app.local
    ```

    If an older `app.local` entry exists, update it so there is only one
    active entry for this hostname.

4.  Test from PowerShell:

    ``` powershell
    curl.exe -k -v https://app.local
    ```

    A successful response should include `HTTP/1.1 200 OK` and the
    application HTML.

**Note:** The WSL IP address can change after WSL or the machine
restarts. If Windows access stops working, check the current WSL IP
using `hostname -I` in Ubuntu and update the `connectaddress` in the
port-proxy rule. The port-forward process must remain running.

### 8. Final verification commands

Run these from Ubuntu:

``` bash
kubectl get nodes
kubectl get pods -n ingress-nginx
kubectl get deployment sprint2-go-deployment
kubectl get service sprint2-go-service
kubectl get secret app-local-tls
kubectl get ingress app-local-ingress
```

Test the HTTPS endpoint using the appropriate `curl` command from Step 6
or Step 7.

## Technical Explanation

### What is Ingress?

Ingress is a Kubernetes API resource that defines rules for routing HTTP
and HTTPS traffic from outside the cluster to Services inside the
cluster. In this project, requests with the hostname `app.local` and
path `/` are routed to `sprint2-go-service`.

### What is the NGINX Ingress Controller?

An Ingress resource defines routing rules, but it needs an Ingress
Controller to implement those rules. The NGINX Ingress Controller
watches Ingress resources and configures NGINX to route incoming
requests.

### What is TLS/HTTPS?

TLS encrypts traffic between the client and server. HTTPS is HTTP
carried over TLS. The Ingress Controller uses the certificate and
private key stored in the Kubernetes TLS Secret to serve HTTPS for
`app.local`.

### What is a Kubernetes TLS Secret?

A TLS Secret stores the certificate and private key in Kubernetes using
the type `kubernetes.io/tls`. The Ingress resource refers to this Secret
using `secretName: app-local-tls`.

### Request flow

1.  The client requests `https://app.local`.
2.  The hostname resolves to the local endpoint configured for the test.
3.  Traffic reaches the NGINX Ingress Controller.
4.  The controller presents the TLS certificate for `app.local`.
5.  The Ingress rule matches the hostname and path.
6.  The controller forwards the request to `sprint2-go-service:3000`.
7.  The Service routes the request to an application Pod.
8.  The Go application returns the response to the client.

## Security Hardening

-   The TLS private key (`tls.key`) must be kept private and must not be
    committed to the repository.
-   Do not store production private keys in source control.
-   A self-signed certificate is used only for local testing. Production
    environments should use a certificate issued by a trusted
    Certificate Authority.
-   The `curl -k` option disables certificate verification. Use it only
    for this self-signed local lab; do not use it as a normal production
    security practice.
-   Kubernetes Secrets are not automatically a complete
    secret-management solution. Production clusters should use suitable
    access controls and secret-management practices.

## Help

### Ingress has no address or does not route traffic

Check the controller and Ingress:

``` bash
kubectl get pods -n ingress-nginx
kubectl get ingress app-local-ingress
kubectl describe ingress app-local-ingress
```

### The application Service is missing

Check the Sprint 2 resources:

``` bash
kubectl get deployment sprint2-go-deployment
kubectl get service sprint2-go-service
kubectl get endpoints sprint2-go-service
```

Ensure the Service has ready endpoints and listens on port `3000`.

### The TLS Secret is missing

Check and recreate the Secret if necessary. Run the creation command
from the directory containing `tls.crt` and `tls.key`:

``` bash
kubectl get secret app-local-tls
kubectl create secret tls app-local-tls --cert=tls.crt --key=tls.key
```

If the Secret already exists, update it with:

``` bash
kubectl create secret tls app-local-tls --cert=tls.crt --key=tls.key --dry-run=client -o yaml | kubectl apply -f -
```

### Windows cannot reach `app.local`

-   Confirm the `kubectl port-forward` process is still running.
-   Confirm the hosts file contains `127.0.0.1 app.local`.
-   Check the current WSL IP using `hostname -I`.
-   Update the Windows port-proxy rule if the WSL IP changed.

## Run the Project

### 1. Go to the Sprint 4 directory

```bash
cd ~/Devops-Intership-2026/sprint-4/ingress
```

### 2. Start Minikube

```bash
minikube start
```

### 3. Check Minikube status

```bash
minikube status
```

### 4. Check the NGINX Ingress Controller

```bash
kubectl get pods -n ingress-nginx
kubectl get service -n ingress-nginx
```

### 5. Check the TLS Secret

```bash
kubectl get secret app-local-tls
```

### 6. Check the Ingress resource

```bash
kubectl get ingress app-local-ingress
```

### 7. Check Kubernetes resources

```bash
kubectl get pods
kubectl get deployment
kubectl get service
```

### 8. Start port forwarding for HTTPS

```bash
kubectl port-forward --address 0.0.0.0 \
  -n ingress-nginx service/ingress-nginx-controller 8443:443
```

Keep this terminal open.

### 9. Test the application

Open another terminal and run:

```bash
curl -k -v --resolve app.local:8443:127.0.0.1 https://app.local:8443/
```

### 10. Test the health endpoint

```bash
curl -k --resolve app.local:8443:127.0.0.1 \
  https://app.local:8443/health
```

Expected output:

```text
OK
```

### 11. Optional: Test from Windows PowerShell

```powershell
curl.exe -k https://app.local
```

Expected result: HTTP `200 OK` and the application HTML.

**Note:** The `-k` option skips certificate verification because this project uses a self-signed TLS certificate. Keep the port-forward terminal running while testing.


## Authors

**Anandhrosh**

DevOps Internship 2026

GitHub: https://github.com/Anandhrosh

## Version History

-   **0.1.0** --- Initial Sprint 4 implementation with NGINX Ingress,
    hostname routing, and self-signed TLS.

## License

This project is intended for educational purposes as part of the DevOps Internship 2026.

## Acknowledgments

-   Kubernetes documentation
-   Minikube documentation
-   NGINX Ingress Controller documentation
