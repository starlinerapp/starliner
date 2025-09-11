# Database Migrations

## Schema generation & migration

To generate the schema required by Better Auth, run the following command:

```
npx @better-auth/cli@latest generate --output ./app/db/schema.ts
```

To generate and apply the migration, run the following commands:

```
npx drizzle-kit generate # generate the migration file
npx drizzle-kit migrate # apply the migration
```
