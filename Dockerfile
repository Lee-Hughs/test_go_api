FROM golang:alpine as build-stage
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.sum ./
RUN export GIN_MODE=release
WORKDIR $GOPATH/src/LeeHughsVariant/go_test_api
COPY . .
RUN mkdir /build
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o /build .

FROM alpine:latest
ENV ITEMS_PATH=/app/items.json
RUN apk update &&\
    apk add ca-certificates
WORKDIR /app
COPY --from=build-stage /build/* /app
COPY items.json .
EXPOSE 8080
ARG UID=1001
RUN adduser -h /app -D -u $UID -g $UID app app && chown -R app:app /app
USER $UID:$UID
CMD ["/app/go_test_api"]