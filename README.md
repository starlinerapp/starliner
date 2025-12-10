# Starliner
The repository is organized into three primary components:
- the **client** folder contains the React Router frontend & backend
- the **server** folder contains the Go backend (Gin) that exposes the REST API.
- the **.docker** folder contains the Docker Files for all supporting services.

## Local Development
You can run the project locally in a few steps using Docker Compose.

### SSL Certificates
Generate a self-signed ssh certificate using the following command:
```bash
openssl req -x509 -newkey rsa:4096 -sha256 -days 365 \
  -nodes -keyout .docker/nginx/ssl/dev.starliner.app.key -out .docker/nginx/ssl/dev.starliner.app.crt \
  -subj "/CN=dev.starliner.app" \
  -addext "subjectAltName=DNS:dev.starliner.app"
```

Then add the development domain to your hosts file:
```bash
grep -qXF "127.0.0.1 dev.starliner.app" /etc/hosts || echo "127.0.0.1 dev.starliner.app" >> /etc/hosts
```

### Docker Compose
1. Navigate to the root directory of the project (where the `compose.yml` file is located).
2. Open a terminal in that directory.
3. Run the following command to build and start the Docker containers: `docker compose up --build -d`

This command will build the necessary Docker images and start the containers defined in the `compose.yml` file.