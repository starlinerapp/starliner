from fastmcp import FastMCP
from pydantic import BaseModel
import asyncio
import logging
import os
import httpx
from src.infrastructure.auth.session import get_user_id_from_token


class IngressPath(BaseModel):
    pathType: str
    path: str
    serviceName: str


class IngressHost(BaseModel):
    host: str
    paths: list[IngressPath]


logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)
mcp = FastMCP("starliner")

AUTH_BASE_URL = os.getenv("AUTH_BASE_URL", "http://client:5173/api/auth")
API_BASE_URL = os.getenv("API_URL", "http://server-api:9090")
BASIC_AUTH_USER = os.getenv("AUTH_USER", "test")
BASIC_AUTH_PASS = os.getenv("AUTH_PASS", "test")

@mcp.tool()
async def login(email: str, password: str):
    """Login to starliner.

    Args:
        email: Email address of the user. Ask the user to input their email.
        password: Password of the user. Ask the user to input their password.
    Returns:
        A bearer token. This token MUST be used in subsequent tool calls
        as the `token` argument.
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{AUTH_BASE_URL}/sign-in/email",
            files={"email": (None, email), "password": (None, password)},
        )

        response.raise_for_status()
        return response.json()["token"]


@mcp.tool()
async def deploy_from_git(
    token: str,
    gitUrl: str,
    environmentId: int,
    serviceName: str,
    port: int,
    projectRepositoryPath: str,
    dockerfilePath: str,
    envs: list[dict[str, str]] = [],
):
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
        token: The bearer token from the login tool call.
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
    user_id = await get_user_id_from_token(token)

    async with httpx.AsyncClient() as client:
        auth = httpx.BasicAuth(BASIC_AUTH_USER, BASIC_AUTH_PASS)
        headers = {"X-User-ID": str(user_id)}
        json = {
            "environmentId": environmentId,
            "serviceName": serviceName,
            "port": port,
            "gitUrl": gitUrl,
            "projectRepositoryPath": projectRepositoryPath,
            "dockerfilePath": dockerfilePath,
            "envs": envs,
        }
        response = await client.post(
            f"{API_BASE_URL}/deployments/git", headers=headers, json=json, auth=auth
        )
        response.raise_for_status()
        return {"status": "ok"}


@mcp.tool()
async def deploy_ingress(
    token: str, environmentId: int, ingressHosts: list[IngressHost]
):
    """Deploy a Traefik ingress to expose HTTP(S) services in an environment.

    Before calling this tool, you MUST:
      1. Call get_environment_deployments(environment_id) to discover existing services.
         - Use the service names from existing deployments to populate the `serviceName` field.
      2. Ask the user for the host prefix(es) and path routing rules if not provided.
      3. Confirm the ingress configuration with the user before deploying.

    Args:
        token: The bearer token from the login tool call.
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
    user_id = await get_user_id_from_token(token)

    async with httpx.AsyncClient() as client:
        auth = httpx.BasicAuth(BASIC_AUTH_USER, BASIC_AUTH_PASS)
        headers = {"X-User-ID": str(user_id)}
        json = {
            "environmentId": environmentId,
            "ingressHosts": [h.model_dump() for h in ingressHosts],
        }
        response = await client.post(
            f"{API_BASE_URL}/deployments/ingresses",
            headers=headers,
            json=json,
            auth=auth,
        )
        response.raise_for_status()
        return {"status": "ok"}


@mcp.tool()
async def get_environments(token: str, project_id: int) -> list[dict]:
    """Get all environments for a project.

    Args:
        token: The bearer token from the login tool call.
        project_id: The ID of the project to get environments for. Prompts the user to
         provide it.

    Returns:
        A list of environments with id, name, and slug.
    """
    user_id = await get_user_id_from_token(token)

    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{API_BASE_URL}/projects/{project_id}/environments",
            headers={"X-User-ID": str(user_id)},
            auth=httpx.BasicAuth(BASIC_AUTH_USER, BASIC_AUTH_PASS),
        )
        response.raise_for_status()
        return response.json()


@mcp.tool()
async def get_environment_deployments(token: str, environment_id: int) -> dict:
    """Get all deployments for an environment.

    Args:
        token: The bearer token from the login tool call.
        environment_id: The ID of the environment to get deployments for.

    Returns:
        A dict containing deployments (git, images, ingresses, databases).
    """
    user_id = await get_user_id_from_token(token)
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{API_BASE_URL}/environments/{environment_id}/deployments",
            headers={"X-User-ID": user_id},
            auth=httpx.BasicAuth(BASIC_AUTH_USER, BASIC_AUTH_PASS),
        )
        response.raise_for_status()
        return response.json()


@mcp.tool()
async def get_organizations(token: str) -> list[dict]:
    """Get all organizations for a user
    Args:
        token: The bearer token from the login tool call.

    Returns:
        A list of organizations with id, name, ownerId, and slug.
    """
    user_id = await get_user_id_from_token(token)
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{API_BASE_URL}/organizations",
            headers={"X-User-ID": str(user_id)},
            auth=httpx.BasicAuth(BASIC_AUTH_USER, BASIC_AUTH_PASS),
        )
        response.raise_for_status()
        return response.json()


@mcp.tool()
async def get_projects(token: str, organization_id: int) -> list[dict]:
    """Get all projects for an organization.

    Args:
        token: The bearer token from the login tool call.
        organization_id: The ID of the organization to get projects for.

    Returns:
        A list of projects with id, name, environments, clusterId, and createdAt.
    """
    user_id = await get_user_id_from_token(token)
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{API_BASE_URL}/organizations/{organization_id}/projects",
            headers={"X-User-ID": str(user_id)},
            auth=httpx.BasicAuth(BASIC_AUTH_USER, BASIC_AUTH_PASS),
        )
        response.raise_for_status()
        return response.json()


@mcp.tool()
async def get_deployment_logs(token: str, deployment_id: int, collect_seconds: int = 5) -> str:
    """Get logs for a specific deployment.

    IMPORTANT: If you don't know the deployment_id, first call get_environment_deployments
    to list available deployments for the environment, then ask the user which deployment
    they want to see logs for.

    Args:
        token: The bearer token from the login tool call.
        deployment_id: The ID of the deployment to get logs for.
        collect_seconds: How many seconds to collect logs for (default 5).

    Returns:
        The deployment logs as a string.
    """
    logger.info(f"Getting logs for deployment {deployment_id} for {collect_seconds}s")
    auth = httpx.BasicAuth(BASIC_AUTH_USER, BASIC_AUTH_PASS)
    logs: list[str] = []

    async with httpx.AsyncClient() as client:
        user_id = await get_user_id_from_token(token)

    async def collect_logs():
        timeout = httpx.Timeout(10.0, read=None)
        async with httpx.AsyncClient(timeout=timeout) as client:
            url = f"{API_BASE_URL}/deployments/{deployment_id}/logs"
            logger.info(f"Connecting to {url}")
            async with client.stream(
                    "GET",
                    url,
                    headers={"X-User-ID": str(user_id)},
                    auth=auth,
            ) as response:
                logger.info(f"Response status: {response.status_code}")
                response.raise_for_status()
                async for line in response.aiter_lines():
                    logger.debug(f"Raw line: {line}")
                    if line.startswith("data:"):
                        log_line = line[5:].strip()
                        logger.info(f"Log: {log_line}")
                        logs.append(log_line)

    try:
        await asyncio.wait_for(collect_logs(), timeout=collect_seconds)
    except asyncio.TimeoutError:
        logger.info(f"Timeout after {collect_seconds}s, collected {len(logs)} lines")
    except httpx.HTTPStatusError as e:
        logger.error(f"HTTP error: {e}")
        return f"Error: {e}"
    except Exception as e:
        logger.error(f"Error: {e}")
        return f"Error: {e}"

    if not logs:
        return "No logs received. The deployment may not be running or producing logs."
    return "\n".join(logs)


app = mcp.http_app()


def create_app():
    return mcp.http_app(transport="streamable-http")


if __name__ == "__main__":
    import uvicorn

    app = create_app()
    uvicorn.run(app, host="0.0.0.0", port=8080)
