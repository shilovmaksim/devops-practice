FROM golang:1.16.7-alpine3.14
WORKDIR /usr/api_server
COPY . .
#RUN go mod tidy && go mod vendor
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod vendor -o api_server/bin/api_server ./api_server/run

FROM alpine:3.14.0
WORKDIR /usr/app
COPY --from=0 /usr/api_server/api_server/bin/api_server .
COPY --from=0 /usr/api_server/api_server/run/config.yml ./run/
CMD ["./api_server"]