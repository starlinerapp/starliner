# Starliner
Repository structure:
- the **client** folder contains the React Router frontend & backend.
- the **server** folder contains the backend services and core platform logic. It is organized into several key modules:
  - **cmd** contains the application entrypoints for the different services.
  - **internal** contains the main application code, organized by service:
    - **api** contains the REST API, through which all client interactions are handled.
    - **builder** is responsible for building Docker images from application source code.
    - **cluster** manages communication and operations with the Kubernetes cluster.
    - **provisioner** handles infrastructure provisioning and configuration.
    - **core** contains shared utilities and cross-cutting functionality used across services.
- the **.docker** folder contains the Docker Files for all supporting services.

## Local Development
You can run the project locally in a few steps using Docker Compose.

### SSL Certificates
Generate a self-signed ssh certificate using the following command and add it to your trust store.
```bash
openssl req -x509 -newkey rsa:4096 -sha256 -days 365 \
  -nodes -keyout .docker/traefik/certs/dev.starliner.app.key -out .docker/traefik/certs/dev.starliner.app.crt \
  -subj "/CN=dev.starliner.app" \
  -addext "subjectAltName=DNS:dev.starliner.app"
```
```bash
openssl req -x509 -newkey rsa:4096 -sha256 -days 365 \
  -nodes -keyout .docker/traefik/certs/auth.dev.starliner.app.key -out .docker/traefik/certs/auth.dev.starliner.app.crt \
  -subj "/CN=auth.dev.starliner.app" \
  -addext "subjectAltName=DNS:auth.dev.starliner.app"
```

Then add the development domain to your hosts file:
```bash
grep -qXF "127.0.0.1 dev.starliner.app auth.dev.starliner.app" /etc/hosts || echo "127.0.0.1 dev.starliner.app auth.dev.starliner.app" >> /etc/hosts
```

### Docker Compose
1. Navigate to the root directory of the project (where the `compose.yml` file is located).
2. Open a terminal in that directory.
3. Run the following command to build and start the Docker containers: `docker compose up --build -d`

This command will build the necessary Docker images and start the containers defined in the `compose.yml` file.

## Contributing
We are not actively accepting external contributions at this stage. While you are welcome to open issues or pull requests, please note that they may be closed, deferred, or left unreviewed.

As the project is still in its early stages, we are currently focusing on maintaining a clear scope, ensuring code quality, and establishing the overall direction of the project. Once the project matures, we plan to revisit and expand opportunities for community contributions.