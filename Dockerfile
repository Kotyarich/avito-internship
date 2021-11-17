FROM golang:1.16.10-alpine as builder
COPY go.mod go.sum /go/src/
WORKDIR /go/src/
RUN go mod download
COPY . /go/src/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/ ./...

FROM alpine
#RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/build /usr/bin/avito-intership
EXPOSE 5555 5555
#RUN chmod +x /usr/bin/avito-intership
ENTRYPOINT ["/usr/bin/avito-intership/balance"]