# Start from golang base image
FROM golang:1.20.3-alpine3.17 as builder

WORKDIR /app

COPY cmd cmd
COPY internal internal
COPY node_modules node_modules
COPY public public
COPY src src
COPY go.mod go.mod
COPY go.sum go.sum
COPY package.json package.json
COPY package-lock.json package-lock.json

RUN apk add --update nodejs npm
RUN npm run build
RUN go build -o app ./cmd

FROM alpine:latest

EXPOSE 8080

WORKDIR /app/

COPY --from=builder /app/app ./
COPY --from=builder /app/build ./build

CMD ["./app"]