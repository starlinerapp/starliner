FROM node:22-alpine AS development-dependencies-env
WORKDIR /app

COPY auth/package*.json ./
RUN npm ci

COPY auth ./

FROM node:22-alpine AS build
WORKDIR /app

COPY --from=development-dependencies-env /app ./
RUN npm run build

FROM node:22-alpine AS production
WORKDIR /app

COPY auth/package*.json ./
RUN npm ci --omit=dev

COPY --from=build /app/dist ./dist

CMD ["npm", "run", "start"]