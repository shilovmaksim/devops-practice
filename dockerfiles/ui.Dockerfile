FROM node:16-alpine3.14
WORKDIR /app
COPY ./ui/app/package-lock.json .
COPY ./ui/app/package.json .
RUN npm ci
COPY ./ui/app .
RUN npm run build

FROM nginx:1.21.1-alpine
COPY ./ui/nginx/nginx.conf /etc/nginx/
COPY --from=0 /app/build/ /usr/share/nginx/html