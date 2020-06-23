FROM golang:1.14.3-alpine3.11 as builder
WORKDIR /go/baseweb/
COPY . .
RUN go build

FROM alpine:3.11
WORKDIR /go/
COPY --from=builder /go/baseweb/baseweb .
EXPOSE 8080
CMD ["./baseweb"]
