FROM alpine:latest

WORKDIR /
COPY . /

EXPOSE 8080
ENTRYPOINT ["./go-gin-example"]