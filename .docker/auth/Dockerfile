FROM node:22-alpine AS development
WORKDIR /app

COPY ../../auth/package*.json ./
RUN npm ci

COPY ../../auth /app
CMD ["npm", "run", "dev"]

FROM node:22-alpine AS build

WORKDIR /app

COPY ../../auth/package.json /app/package.json
COPY ../../auth/package-lock.json /app/package-lock.json
RUN npm ci

COPY ../../auth /app

RUN npm run build

FROM node:22-alpine

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000

COPY --from=build /app/package.json /app/package.json
COPY --from=build /app/package-lock.json /app/package-lock.json
RUN npm ci --omit=dev

COPY --from=build /app/dist /app/dist

RUN chown -R node:node /app
USER node

EXPOSE 3000
CMD ["node", "dist/index.js"]