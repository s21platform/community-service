FROM golang:1.24-alpine AS builder

WORKDIR /usr/src/service
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build -o build/main cmd/service/main.go
RUN go build -o build/workerParticipant cmd/worker/patricipants/main.go
RUN go build -o build/workerCampus cmd/worker/campuses/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /usr/src/service/build/main .
COPY --from=builder /usr/src/service/build/workerParticipant .
COPY --from=builder /usr/src/service/build/workerCampus .
CMD ./main & ./workerParticipant & ./workerCampus