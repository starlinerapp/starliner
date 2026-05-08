import httpx

from src.utils.env import get_required_env

AUTH_SESSION_URL = get_required_env("AUTH_BASE_URL") + "/get-session"


class UnauthorizedError(Exception):
    pass


async def get_user_id_from_token(token: str) -> int:
    async with httpx.AsyncClient() as client:
        response = await client.get(
            AUTH_SESSION_URL,
            headers={"Authorization": f"Bearer {token}"},
        )

        if response.status_code != 200:
            raise UnauthorizedError("Invalid or expired session token")

        data = response.json()

        user = data.get("user")
        if not user:
            raise UnauthorizedError("No user in session")

        user_id = user.get("id")
        if not user_id:
            raise UnauthorizedError("User ID missing in session")

        return user_id
