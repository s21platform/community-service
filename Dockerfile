FROM golang:1.22-alpine AS builder

WORKDIR /usr/src/service
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o build/main cmd/service/main.go
RUN go build -o build/workerLogins cmd/worker/logins/main.go
RUN go build -o build/workerCampus cmd/worker/campuses/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /usr/src/service/build/main .
COPY --from=builder /usr/src/service/build/workerLogins .
COPY --from=builder /usr/src/service/build/workerCampus .
CMD ./main & ./workerLogins & ./workerCampus