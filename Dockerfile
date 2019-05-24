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
RUN apk add --no-cache openssl
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz
COPY --from=builder /main ./
RUN chmod +x ./main
RUN apk update && apk add bash
CMD dockerize -wait tcp://db:3306  -timeout 15s && ./main
EXPOSE 2004