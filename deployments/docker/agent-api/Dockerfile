FROM golang:-alpine as builder

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"

COPY .  .
RUN export GO111MODULE=on && CGO_ENABLED=0 GOOS=linux go build -o build/agent-api agent-api/main.go


FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk --no-cache add ca-certificates


RUN mkdir /app
WORKDIR /app

COPY --from=builder /root/build/agent-api ./

ENTRYPOINT ["/app/agent-api"]
