FROM golang:1.21-alpine AS builder

# RUN apk update && apk add --no-cache git
RUN apk update

WORKDIR /app
COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -o /go/bin/sdn .


FROM alpine:3.19

ENV DATABASE_URL=postgresql://localhost:4321/postgres

WORKDIR /app

COPY --from=builder /go/bin/sdn /app/sdn

CMD ["/bin/sh", "-c", "/app/sdn"]