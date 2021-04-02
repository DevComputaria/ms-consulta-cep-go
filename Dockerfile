FROM golang:1.16 as builder
WORKDIR /go/src/devcomputaria
COPY . .
RUN GOOS=linux go build -ldflags="-s -w"

FROM scratch
WORKDIR /go/src/devcomputaria
COPY --from=builder /go/src/devcomputaria/devcomputaria .

CMD ["./devcomputaria"]

