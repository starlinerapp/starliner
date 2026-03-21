from fastmcp import FastMCP
from pydantic import BaseModel
import os
import httpx


class IngressPath(BaseModel):
    pathType: str
    path: str
    serviceName: str


class IngressHost(BaseModel):
    host: str
    paths: list[IngressPath]

mcp = FastMCP("starliner")

API_BASE_URL = os.getenv("STARLINER_API_URL", "http://server-api:9090")
BASIC_AUTH_USER = os.getenv("STARLINER_AUTH_USER", "test")
BASIC_AUTH_PASS = os.getenv("STARLINER_AUTH_PASS", "test")
USER_ID = os.getenv("STARLINER_USER_ID", "RfDNq6DidOFZzk1rrBzkM73vCUhLyPgH")


@mcp.tool()
async def hello() -> str:
    """Returns a hello world message."""
    return "Hello World"

@mcp.tool()
async def deploy_from_git(gitUrl: str, environmentId: int, serviceName: str, port: int, projectRepositoryPath: str, dockerfilePath: str, envs: list[dict[str, str]] = []):
    """Deploy a service from a Git repository.

    Before calling this tool, you MUST:
      1. Ask the user which project and environment they want to deploy to.
         - Use get_projects(organization_id) to list available projects.
         - Use get_environments(project_id) to list available environments for the chosen project.
      2. Call get_environment_deployments(environment_id) to discover existing deployments in the target environment.
         - Use credentials and connection details from existing deployments (e.g. databases) to automatically populate the env vars the application needs.
      3. Infer git_url, dockerfile_path, and project_repository_path from the repo.
      4. Try to create a Dockerfile if there is none.
      5. Infer port from the Dockerfile (e.g. EXPOSE directive).
      6. Scan the application code for any additional required env var names and resolve their values from step 2 or ask the user.
      7. Ask the user for: service_name and any env var values you could not resolve automatically.

    Args:
        gitUrl: The Git repository URL.
        environmentId: The ID of the target environment.
        serviceName: The name of the service to deploy.
        port: The port the application listens on.
        projectRepositoryPath: Root directory of the project in the repo.
        dockerfilePath: Path to the Dockerfile in the repo, relative to the projectRepositoryPath.
        envs: List of environment variables, each as {"name": "VAR_NAME", "value": "VAR_VALUE"}.

    Returns:
        The deployment response from the backend.
    """
    async with httpx.AsyncClient() as client:
        auth = httpx.BasicAuth(BASIC_AUTH_USER, BASIC_AUTH_PASS)
        headers = {"X-User-ID": USER_ID}
        json = {"environmentId": environmentId, "serviceName": serviceName, "port": port, "gitUrl": gitUrl, "projectRepositoryPath": projectRepositoryPath, "dockerfilePath": dockerfilePath, "envs": envs}
        response = await client.post(f"{API_BASE_URL}/deployments/git", headers=headers, json=json, auth=auth)
        response.raise_for_status()
        return {"status": "ok"}


@mcp.tool()
async def deploy_ingress(environmentId: int, ingressHosts: list[IngressHost]):
    """Deploy a Traefik ingress to expose HTTP(S) services in an environment.

    Before calling this tool, you MUST:
      1. Call get_environment_deployments(environment_id) to discover existing services.
         - Use the service names from existing deployments to populate the `serviceName` field.
      2. Ask the user for the host prefix(es) and path routing rules if not provided.
      3. Confirm the ingress configuration with the user before deploying.

    Args:
        environmentId: The ID of the target environment.
        ingressHosts: List of host routing rules. Each entry has the shape:
            {
                "host": str,        # Subdomain prefix (e.g. "myapp" → "myapp.<cluster-domain>")
                "paths": [
                    {
                        "pathType": str,    # "Prefix" or "Exact"
                        "path": str,        # URL path to match, e.g. "/" or "/api"
                        "serviceName": str  # Name of the target service to route traffic to
                    }
                ]
            }

    Returns:
        The deployment response from the backend.
    """
    async with httpx.AsyncClient() as client:
        auth = httpx.BasicAuth(BASIC_AUTH_USER, BASIC_AUTH_PASS)
        headers = {"X-User-ID": USER_ID}
        json = {"environmentId": environmentId, "ingressHosts": [h.model_dump() for h in ingressHosts]}
        response = await client.post(f"{API_BASE_URL}/deployments/ingresses", headers=headers, json=json, auth=auth)
        response.raise_for_status()
        return {"status": "ok"}




@mcp.tool()
async def get_environments(project_id: int) -> list[dict]:
    """Get all environments for a project.

    Args:
        project_id: The ID of the project to get environments for.

    Returns:
        A list of environments with id, name, and slug.
    """
    auth = httpx.BasicAuth(BASIC_AUTH_USER, BASIC_AUTH_PASS)
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{API_BASE_URL}/projects/{project_id}/environments",
            headers={"X-User-ID": USER_ID},
            auth=auth,
        )
        response.raise_for_status()
        return response.json()


@mcp.tool()
async def get_environment_deployments(environment_id: int) -> dict:
    """Get all deployments for an environment.

    Args:
        environment_id: The ID of the environment to get deployments for.

    Returns:
        A dict containing deployments (git, images, ingresses, databases).
    """
    auth = httpx.BasicAuth(BASIC_AUTH_USER, BASIC_AUTH_PASS)
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{API_BASE_URL}/environments/{environment_id}/deployments",
            headers={"X-User-ID": USER_ID},
            auth=auth,
        )
        response.raise_for_status()
        return response.json()


@mcp.tool()
async def get_projects(organization_id: int) -> list[dict]:
    """Get all projects for an organization.

    Args:
        organization_id: The ID of the organization to get projects for.

    Returns:
        A list of projects with id, name, environments, clusterId, and createdAt.
    """
    auth = httpx.BasicAuth(BASIC_AUTH_USER, BASIC_AUTH_PASS)
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{API_BASE_URL}/organizations/{organization_id}/projects",
            headers={"X-User-ID": USER_ID},
            auth=auth,
        )
        response.raise_for_status()
        return response.json()


app = mcp.http_app()
