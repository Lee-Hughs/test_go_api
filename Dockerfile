FROM golang:alpine
ENV CGO_ENABLED=0
ENV ITEMS_PATH=./items.json
WORKDIR /app
COPY go.mod go.sum items.json ./
RUN export GIN_MODE=release
WORKDIR $GOPATH/src/LeeHughsVariant/go_test_api
COPY . .
RUN mkdir /build
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o /build .
RUN apk update &&\
    apk add ca-certificates

WORKDIR /app
RUN mv /build/* /app

EXPOSE 8080

ARG UID=1001
RUN adduser -h /app -D -u $UID -g $UID app app && chown -R app:app /app
USER $UID:$UID
CMD ["/app/go_test_api"]