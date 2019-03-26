# Build stage
FROM golang:alpine as build
ARG VERSION
RUN apk --update add --no-cache make git

WORKDIR /go/src/app
COPY . .
RUN make build && mv server /usr/local/bin

# Final stage
FROM alpine
COPY --from=build /usr/local/bin/server /usr/local/bin
CMD ["/usr/local/bin/server"]
