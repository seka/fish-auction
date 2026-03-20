FROM node:20-alpine

WORKDIR /app

COPY package.json yarn.lock panda.config.ts ./
RUN yarn install --frozen-lockfile --ignore-scripts

COPY . .
RUN yarn prepare

EXPOSE 3000

CMD ["yarn", "dev"]
