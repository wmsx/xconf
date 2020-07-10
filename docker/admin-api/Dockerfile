FROM golang:alpine as builder

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"

COPY .  .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/admin-api admin-api/main.go


FROM node:12-alpine as node-builder

RUN mkdir /app
WORKDIR /app

COPY . .
RUN npm rebuild node-sass && npm install && npm run build

FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/build/admin-api ./
COPY --from=node-builder /app/build dist

ENTRYPOINT ["/app/admin-api"]
