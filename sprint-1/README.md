# Sprint 1: Containerization & Image Optimization

## Project Title

**Optimized Multi-Stage Dockerfile for a Go Application**

## Description

This project is part of the DevOps Internship 2026 evaluation. The objective is to build an optimized, secure, and lightweight Docker image for a sample Go application using a multi-stage Dockerfile.

The project separates the Go build environment from the production runtime image. The application is compiled in a Go Alpine builder stage, and only the compiled executable is copied into a minimal `scratch` runtime image. The resulting image is approximately **6.02 MB**, which is below the project target of 100 MB.

Security is a key focus of this project. The application runs under a custom non-root user named `appuser`. The Dockerfile also configures a health check, exposes port 3000, and uses a `.dockerignore` file to keep unnecessary files out of the build context.

### Project Objectives

* Implement a multi-stage Docker build.
* Optimize the final Docker image to less than 100 MB.
* Use a lightweight build image and a minimal `scratch` runtime image.
* Run the application using a custom non-root user.
* Compile the Go application as a static Linux executable.
* Configure a Docker health check.
* Expose the application on port 3000.
* Document the build, execution, and verification procedures.

## Getting Started

### Dependencies

The following tools and software are required:

* Ubuntu 24.04 LTS (or another supported Linux environment)
* Docker Engine
* Git
* GitHub account

The Go compiler is provided by the Docker build stage, so Go does not need to be installed separately on the host.

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

The directory should contain the following project files:

```text
sprint-1/
├── .dockerignore
├── Dockerfile
├── main.go
└── README.md
```

## Executing Program

**Step 1: Build the Docker image**

Build the application image using the Dockerfile:

```bash
docker build -t sprint1-go-app .
```

**Step 2: Verify the Docker image**

```bash
docker images sprint1-go-app
```

**Step 3: Check the image size**

```bash
docker image inspect sprint1-go-app --format '{{.Size}}'
```

The final image should be less than 100 MB. In the tested build, the image size was approximately **6.02 MB**. Docker may display the size in bytes in the inspect command.

**Step 4: Run the Docker container**

Run the application and map host port 3000 to container port 3000:

```bash
docker run -d \
  --name sprint1-go-container \
  -p 3000:3000 \
  sprint1-go-app
```

If host port 3000 is already in use, use another host port, such as 8080:

```bash
docker run -d \
  --name sprint1-go-container \
  -p 8080:3000 \
  sprint1-go-app
```

When using port 8080, access the application at `http://localhost:8080`.

**Step 5: Verify the running container**

```bash
docker ps
```

Check that the container is running and that the port mapping is shown as `3000->3000` (or `8080->3000` if you used the alternative command).

**Step 6: Test the application**

Access the application using curl:

```bash
curl http://localhost:3000
```

Expected output:

```text
Welcome to Sprint 1 Docker Project!
```

You can also open `http://localhost:3000` in a web browser.

**Step 7: Test the health endpoint**

```bash
curl http://localhost:3000/health
```

Expected output:

```text
OK
```

**Step 8: Verify non-root execution**

The final `scratch` image does not include common shell utilities such as `whoami` or `id`. Verify the configured container user using Docker inspect:

```bash
docker inspect sprint1-go-container --format='User: {{.Config.User}}'
```

Expected output:

```text
User: appuser
```

To inspect the custom user's account entry in the image, export the container filesystem and read its passwd file:

```bash
docker export sprint1-go-container | tar -xOf - etc/passwd
```

The output should include an entry for `appuser` with UID `10001`. A nonzero UID means the application is not configured to run as root.

**Step 9: Verify the Docker health check**

```bash
docker inspect sprint1-go-container --format '{{.State.Health.Status}}'
```

Expected output after the health check succeeds:

```text
healthy
```

**Step 10: View application logs**

```bash
docker logs sprint1-go-container
```

**Step 11: Stop and remove the container**

```bash
docker stop sprint1-go-container
docker rm sprint1-go-container
```

To start the existing container again later, use:

```bash
docker start sprint1-go-container
```

## Dockerfile Optimization

The Dockerfile uses a multi-stage build to separate the build environment from the production runtime.

### Build Stage

* Uses `golang:1.24-alpine` as the builder image.
* Copies `main.go` into the build stage.
* Compiles the Go application for Linux with CGO disabled.
* Uses `-trimpath` and linker flags `-s -w` to reduce build metadata and executable size.
* Creates a passwd entry for the custom runtime user.

### Production Stage

* Uses `scratch`, an empty base image, as the runtime image.
* Copies only the compiled application executable and the passwd file from the builder stage.
* Runs the application as the custom non-root user `appuser`.
* Exposes port 3000.
* Configures a Docker health check using the application's `healthcheck` command.
* Starts the compiled application executable.

### Image Optimization Techniques

* Multi-stage Docker build to keep the compiler and build tools out of the final image.
* `scratch` runtime image to avoid including an operating system userland or unnecessary packages.
* Static Go compilation with `CGO_ENABLED=0`.
* `-trimpath` and `-ldflags="-s -w"` to reduce executable metadata and size.
* `.dockerignore` to exclude unnecessary files from the build context.

## Security Hardening

The following security measures are implemented:

1. **Non-root execution:** The application runs as the custom `appuser` rather than root.
2. **Minimal runtime image:** The `scratch` image contains only the files required to run the application.
3. **Multi-stage build:** The final image excludes the Go compiler and build environment.
4. **Reduced executable metadata:** Build flags remove path and symbol/debug information that is not needed at runtime.
5. **Docker health check:** The application health endpoint is checked to report container health.
6. **Build context filtering:** `.dockerignore` excludes unnecessary files from the build context.

## Help

### Common Problems and Solutions

**1. Docker image build fails**

Check that the required project files exist and review the build output:

```bash
ls -la
docker build -t sprint1-go-app .
```

**2. Port 3000 is already in use**

Use another host port while keeping the application port inside the container at 3000:

```bash
docker run -d \
  --name sprint1-go-container \
  -p 8080:3000 \
  sprint1-go-app
```

Access the application at `http://localhost:8080`.

**3. Container is not running**

Check the container status:

```bash
docker ps -a
```

View the container logs:

```bash
docker logs sprint1-go-container
```

**4. Container is not healthy**

Check the health status and recent health-check results:

```bash
docker inspect sprint1-go-container \
  --format '{{json .State.Health}}'
```

Check the application logs:

```bash
docker logs sprint1-go-container
```

**5. Verify the container is not running as root**

The `scratch` image does not include `id` or `whoami`. Check the configured user with:

```bash
docker inspect sprint1-go-container --format='User: {{.Config.User}}'
```

The expected user is `appuser`. The Dockerfile assigns this account UID `10001`, which is nonzero.

**6. Docker image exceeds 100 MB**

Check the image size:

```bash
docker images sprint1-go-app
```

Review the final Dockerfile stage and confirm that only the compiled executable and required passwd file are copied into the `scratch` image. Rebuild and verify the image size after making changes.

## Authors

**Anandhrosh**  
DevOps Internship 2026

GitHub: [@Anandhrosh](https://github.com/Anandhrosh)

## Version History

* 1.0
  * Created a sample Go HTTP application.
  * Implemented a multi-stage Dockerfile.
  * Used a minimal `scratch` runtime image.
  * Configured a custom non-root user.
  * Added a Docker health check.
  * Optimized the compiled executable and Docker build context.
  * Added build, execution, and verification instructions.

  * 2.0
  Updated Application Features

* The Go application now includes a simple web interface with a home page and a health-status page.
* Home Page (/)
* Displays a welcome message.
* Provides a simple, user-friendly web interface.
* Runs on port 3000.
* Health Page (/health)
* Displays the application's health information.
* Returns HTTP status 200 when the application is running.
* Used by Docker's built-in HEALTHCHECK to monitor the application.

## License

This project is intended for educational purposes as part of the DevOps Internship 2026.

## Acknowledgments

* [Docker Documentation](https://docs.docker.com/)
* [Go Documentation](https://go.dev/doc/)
* [Alpine Linux](https://www.alpinelinux.org/)
* [Git Documentation](https://git-scm.com/doc)
* [GitHub Documentation](https://docs.github.com/)
