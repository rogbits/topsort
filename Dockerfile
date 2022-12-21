FROM golang:bullseye AS builder
COPY /lib /app/lib
COPY main.go /app/
COPY go.mod /app/
WORKDIR /app
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o topsort .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/ ./
CMD ["./topsort"]