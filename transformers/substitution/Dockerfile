FROM golang:1.19-alpine as builder
ENV CGO_ENABLED=0
WORKDIR /go/src/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags '-w -s' -v -o /usr/local/bin/function ./

FROM ubuntu:latest as user
RUN useradd -u 10001 scratchuser

FROM scratch
COPY --from=builder /usr/local/bin/function /substition
COPY --from=user /etc/passwd /etc/passwd
ENTRYPOINT ["/substition"]
