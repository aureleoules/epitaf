FROM node:slim as builder

WORKDIR /tmp/app

COPY . .

RUN cp /tmp/app/.env /tmp/app/.env

RUN yarn
RUN yarn build

FROM nginx:stable-alpine
COPY nginx.config /etc/nginx/conf.d/default.conf
COPY --from=builder /tmp/app/build /usr/share/nginx/html

EXPOSE 80