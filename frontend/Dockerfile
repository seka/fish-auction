FROM node:20-alpine

WORKDIR /app

COPY package.json yarn.lock panda.config.ts ./
RUN yarn install --frozen-lockfile

COPY . .

EXPOSE 3000

CMD ["yarn", "dev"]
