FROM node:22-alpine AS development-dependencies-env
COPY ../../client /app
WORKDIR /app
RUN npm ci

FROM node:22-alpine
COPY ../../client /app/
COPY --from=development-dependencies-env /app/node_modules /app/node_modules
WORKDIR /app
RUN npm run build

CMD ["npm", "run", "start"]