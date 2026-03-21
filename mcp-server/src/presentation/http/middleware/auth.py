import httpx

async def _send_401(send, scope):
    if scope["type"] == "http":
        await send({
            "type": "http.response.start",
            "status": 401,
            "headers": [[b"content-type", b"application/json"]],
        })
        await send({
            "type": "http.response.body",
            "body": b'{"error": "Unauthorized: invalid or missing session token"}',
        })

class AuthMiddleware:
    def __init__(self, app):
        self.app = app

    async def __call__(self, scope, receive, send):
        if scope["type"] != "http":
            await self.app(scope, receive, send)
            return

        headers = dict(scope.get("headers", []))
        auth_header = headers.get(b"authorization", b"").decode()

        if not auth_header.startswith("Bearer "):
            await _send_401(send, scope)
            return

        token = auth_header[7:]

        try:
            async with httpx.AsyncClient(timeout=5.0) as client:
                response = await client.get(
                    "http://client:5173/api/auth/get-session", headers={"Authorization": f"Bearer {token}"},
                )
            if response.status_code != 200:
                await _send_401(send, scope)
                return

            data = response.json()

            if not data.get("session") or not data.get("user"):
                await _send_401(send, scope)
                return

        except (httpx.HTTPError, ValueError):
            await _send_401(send, scope)
            return

        scope["user"] = data["user"]
        await self.app(scope, receive, send)