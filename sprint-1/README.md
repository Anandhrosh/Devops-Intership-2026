# Sprint 1: Containerization & Image Optimization

## Project Title

**Optimized Multi-Stage Dockerfile for a Node.js Application**

## Description

This project is part of the DevOps Internship 2026 evaluation. The objective is to build an optimized, secure, and lightweight Docker image for a sample Node.js application using a multi-stage Dockerfile.

The project demonstrates how to separate build dependencies from the production runtime environment to reduce the final Docker image size. It uses Node.js with Alpine Linux as the base image to achieve a target image size of less than 100 MB.

Security is a key focus of this project. The application runs under a custom non-root user instead of root. Docker layer caching is optimized by copying dependency files before application source code. A Docker health check is also configured to monitor the application's availability.

### Project Objectives

* Implement a multi-stage Docker build.
* Optimize the final Docker image to less than 100 MB.
* Use a lightweight Alpine Linux base image.
* Run the application using a custom non-root user.
* Optimize Docker layer caching.
* Configure a Docker health check.
* Expose the application on port 3000.
* Document the build, execution, and verification procedures.

## Getting Started

### Dependencies

The following tools and software are required:

* Ubuntu 24.04 LTS
* Docker Engine
* Node.js 22
* npm
* Git
* GitHub account

### Installing

**1. Clone the repository**

```bash
git clone https://github.com/Anandhrosh/Devops-Intership-2026.git
```

**2. Navigate to the Sprint 1 directory**

```bash
cd Devops-Intership-2026/sprint-1
```

**3. Verify the project files**

```bash
ls -la
```

The directory should contain the following files:

```text
sprint-1/
├── app.js
├── package.json
├── package-lock.json
├── Dockerfile
├── .dockerignore
└── README.md
```

**4. Generate the package lock file if it is missing**

```bash
npm install --package-lock-only
```

### Executing Program

**Step 1: Build the Docker image**

Build the application image using the Dockerfile:

```bash
docker build -t sprint1-node-app:1.0 .
```

**Step 2: Verify the Docker image**

```bash
docker images sprint1-node-app
```

**Step 3: Check the image size**

```bash
docker image inspect sprint1-node-app:1.0 --format '{{.Size}}'
```

The final Docker image must be less than 100 MB to meet the project requirement.

**Step 4: Run the Docker container**

```bash
docker run -d \
  --name sprint1-app \
  -p 3000:3000 \
  sprint1-node-app:1.0
```

**Step 5: Verify the running container**

```bash
docker ps
```

**Step 6: Test the application**

Access the application using curl:

```bash
curl http://localhost:3000
```

Expected output:

```text
Hello from Sprint 1 Docker Project!
```

**Step 7: Test the health endpoint**

```bash
curl http://localhost:3000/health
```

Expected output:

```text
OK
```

**Step 8: Verify non-root execution**

Check the user running inside the container:

```bash
docker exec sprint1-app whoami
```

Expected output:

```text
appuser
```

Verify the user ID:

```bash
docker exec sprint1-app id
```

The container should run with a nonzero UID rather than root (UID 0).

**Step 9: Verify the Docker health check**

```bash
docker inspect sprint1-app --format '{{.State.Health.Status}}'
```

Expected output after a successful health check:

```text
healthy
```

**Step 10: View application logs**

```bash
docker logs sprint1-app
```

**Step 11: Stop and remove the container**

```bash
docker stop sprint1-app
docker rm sprint1-app
```

## Dockerfile Optimization

The Dockerfile uses a multi-stage build to separate the build environment from the production runtime.

### Build Stage

* Uses Node.js Alpine as the base image.
* Copies `package.json` and `package-lock.json` first.
* Installs production dependencies using `npm ci --omit=dev`.
* Copies the application source code.

### Production Stage

* Uses Node.js Alpine as the runtime image.
* Creates a custom non-root user named `appuser`.
* Copies only the required application files and dependencies from the build stage.
* Exposes port 3000.
* Configures a Docker health check.
* Runs the application as the non-root user.

### Image Optimization Techniques

* Multi-stage Docker build to avoid unnecessary build tools in the final image.
* Alpine Linux to reduce the base image size.
* Efficient layer caching by copying dependency manifests before application source code.
* Production-only dependency installation.
* `.dockerignore` to exclude unnecessary files from the build context.

## Security Hardening

The following security measures are implemented:

1. **Non-root execution:** The application runs as the custom `appuser` instead of root.
2. **Minimal runtime image:** Alpine Linux is used to reduce unnecessary packages.
3. **Multi-stage build:** The final image excludes the separate build environment.
4. **Production dependencies:** Only production dependencies are installed.
5. **Docker health check:** The application health endpoint is monitored.
6. **Build context filtering:** `.dockerignore` excludes unnecessary files and local dependencies.

## Help

### Common Problems and Solutions

**1. Docker image build fails**

Check the Dockerfile and ensure all required application files exist.

```bash
ls -la
docker build -t sprint1-node-app:1.0 .
```

**2. Port 3000 is already in use**

Use another host port:

```bash
docker run -d \
  --name sprint1-app \
  -p 8080:3000 \
  sprint1-node-app:1.0
```

Access the application at `http://localhost:8080`.

**3. Container is not running**

Check the container status:

```bash
docker ps -a
```

View the container logs:

```bash
docker logs sprint1-app
```

**4. Container is not healthy**

Check the health status and recent health-check results:

```bash
docker inspect sprint1-app \
  --format '{{json .State.Health}}'
```

Check the application logs:

```bash
docker logs sprint1-app
```

**5. Verify the container is not running as root**

```bash
docker exec sprint1-app id
```

The UID should be nonzero.

**6. Docker image exceeds 100 MB**

Check the image size:

```bash
docker images sprint1-node-app
```

Review the base image and the files copied into the final stage. Rebuild and verify the image size after making optimizations.

## Authors

**Anandhrosh**
DevOps Internship 2026

GitHub: [@Anandhrosh](https://github.com/Anandhrosh)

## Version History

* 1.0

  * Initial Sprint 1 implementation.
  * Created a sample Node.js application.
  * Implemented a multi-stage Dockerfile.
  * Configured a non-root user.
  * Added a Docker health check.
  * Optimized Docker layer caching.
  * Added build, execution, and verification instructions.

## License

This project is intended for educational purposes as part of the DevOps Internship 2026.

## Acknowledgments

* [Docker Documentation](https://docs.docker.com/)
* [Node.js Documentation](https://nodejs.org/docs/)
* [Alpine Linux](https://www.alpinelinux.org/)
* [Git Documentation](https://git-scm.com/doc)
* [GitHub Documentation](https://docs.github.com/)
