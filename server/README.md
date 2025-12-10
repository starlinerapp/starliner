# Starliner Server

## Setting up Delve for debugging
1. In Goland click on `"Add Configuration" > "Add new..." > Go Remote`
1. Provide a name (e.g. server) and change the port to `40000`


## Commands
The project uses **sqlc** to generate type-safe database queries and **Goose** to manage and apply database schema migrations.

### sqlc
- Generate Go code from SQL files: 
```bash
sqlc generate
```
- Verify configuration with
```bash
sqlc vet
```
### goose
- Run database migrations:
```bash
goose up
```
- Roll back last migration:
```bash
goose down
```