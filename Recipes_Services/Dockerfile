FROM golang:1.19 AS build
ADD . /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /src/bin/main /usr/local/bin/recipes_app
RUN chmod +x /usr/local/bin/recipes_app
CMD ["/usr/local/bin/recipes_app"]
