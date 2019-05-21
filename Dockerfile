FROM golang:latest AS builder
RUN curl https://glide.sh/get | sh
COPY glide.yaml glide.yaml
COPY glide.lock glide.lock
RUN glide install
ADD . /go/src/wallet
WORKDIR /go/src/wallet
RUN go build && go install
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /main .
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /main ./
RUN chmod +x ./main
ENTRYPOINT ["./main"]
EXPOSE 3030