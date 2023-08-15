FROM golang:1.19 AS build
ADD . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main cmd/server/main.go

FROM alpine:latest
COPY --from=build /src/bin/main /usr/local/bin/rss_worker
RUN chmod +x /usr/local/bin/rss_worker
CMD ["/usr/local/bin/rss_worker"]
