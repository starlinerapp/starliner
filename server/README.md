# Starliner Server

## Setting up Delve for debugging
1. In Goland click on `"Add Configuration" > "Add new..." > Go Remote`
1. Provide a name (e.g., api) and change the port to `40000`

Repeat the above steps for each service you want to debug, changing the port number for each one (e.g., `40001` for builder, `40002` for cluster and `40003` for provisioner).

## Debugging GitHub Webhooks locally
In order to develop your app locally, you can use a webhook proxy URL to forward webhooks from GitHub to your computer.
1. In your browser, navigate to https://smee.io/. 
1. Click Start a new channel.
1. Copy the full URL under `Webhook Proxy URL`. You will use this URL in a later step.
1. On GitHub,
   1. Navigate to `Organization Settings` > `Developer Settings` > `GitHub Apps` > `dev.starliner.app`.
   1. Under "Webhook URL", enter your webhook proxy URL.
1. Open a terminal and run the following command to start forwarding webhooks to your local machine:
```bash
node --use-system-ca $(which smee) -u https://smee.io/[...] -t https://dev.starliner.app/webhooks/github
```

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