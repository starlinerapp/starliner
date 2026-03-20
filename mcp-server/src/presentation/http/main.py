from fastmcp import FastMCP

auth = JWTVerifier(
    jwks_uri=os.environ["BETTER_AUTH_JWKS_URI"],
    issuer=os.environ["BETTER_AUTH_ISSUER"],
    audience=os.environ["BETTER_AUTH_AUDIENCE"],
)

mcp = FastMCP("starliner")


@mcp.tool()
async def hello() -> str:
    """Returns a hello world message."""
    return "Hello World"

