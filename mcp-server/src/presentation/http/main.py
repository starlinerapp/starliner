import logging
import os

import httpx
from fastmcp import FastMCP
from src.infrastructure.auth.session import get_user_id_from_token

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
            files={
                "email": (None, email),
                "password": (None, password)
            },
        )

        response.raise_for_status()
        return response.json()["token"]

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

def create_app():
    return mcp.http_app(transport="streamable-http")

if __name__ == "__main__":
    import uvicorn
    app = create_app()
    uvicorn.run(app, host="0.0.0.0", port=8080)