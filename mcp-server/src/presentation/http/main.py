from fastmcp import FastMCP

mcp = FastMCP("starliner")


@mcp.tool()
async def hello() -> str:
    """Returns a hello world message."""
    return "Hello World"

