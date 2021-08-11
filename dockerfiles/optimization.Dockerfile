FROM golang:1.16-alpine3.14
WORKDIR /usr/optimization_server
COPY . .
#RUN go mod tidy && go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -o ./optimization_server/bin/optimization_server ./optimization_server/run

FROM python:3.9.6-alpine3.14
WORKDIR /usr/optimization_server
COPY --from=0 /usr/optimization_server/optimization_server/bin/ .
COPY --from=0 /usr/optimization_server/optimization_server/python_script/main.py ./python_script/
COPY --from=0 /usr/optimization_server/optimization_server/run/config.yml ./run/
CMD ["./optimization_server"]