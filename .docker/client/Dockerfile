FROM node:22-alpine AS development-dependencies-env
WORKDIR /app
COPY ../../client/package.json ../../client/package-lock.json ./
RUN npm ci
COPY ../../client .

FROM node:22-alpine AS production-dependencies-env
WORKDIR /app
COPY ../../client/package.json ../../client/package-lock.json ./
RUN npm ci --omit=dev

FROM node:22-alpine AS build-env
WORKDIR /app
COPY ../../client .
COPY --from=development-dependencies-env /app/node_modules /app/node_modules
RUN npm run build

FROM node:22-alpine
WORKDIR /app
COPY ../../client/package.json ../../client/package-lock.json ./
COPY --from=production-dependencies-env /app/node_modules /app/node_modules
COPY --from=build-env /app/build /app/build
COPY --from=build-env /app/server.mjs /app/server.mjs
COPY --from=build-env /app/instrument.server.mjs /app/instrument.server.mjs
CMD ["npm", "run", "start"]