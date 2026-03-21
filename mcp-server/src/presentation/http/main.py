import os

import httpx
from fastmcp import FastMCP

mcp = FastMCP("starliner")

API_BASE_URL = os.getenv("STARLINER_API_URL", "http://server-api:9090")
BASIC_AUTH_USER = os.getenv("STARLINER_AUTH_USER", "test")
BASIC_AUTH_PASS = os.getenv("STARLINER_AUTH_PASS", "test")
USER_ID = os.getenv("STARLINER_USER_ID", "")


@mcp.tool()
async def hello() -> str:
    """Returns a hello world message."""
    return "Hello World"


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

