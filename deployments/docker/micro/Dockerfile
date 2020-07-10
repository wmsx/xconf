FROM golang:alpine as builder

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"

COPY .  .
RUN CGO_ENABLED=0 GOOS=linux go build -o build/micro micro/main.go

FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/build/micro ./

EXPOSE 8090
ENTRYPOINT ["/app/micro", "api", "--handler=http"]
