FROM node:22-alpine AS base
WORKDIR /app
COPY package*.json ./
RUN npm install

FROM base AS development 
COPY . .
EXPOSE 3000
CMD ["npm", "run", "dev"]

FROM base AS build
COPY . . 
RUN npm run build

FROM node:22-alpine AS production 
WORKDIR /app
COPY --from=base /app/package*.json ./
COPY --from=base /app/node_modules ./node_modules
COPY --from=build /app/.next ./.next
COPY --from=build /app/public ./public
EXPOSE 3000
CMD ["npm", "run", "start"]





