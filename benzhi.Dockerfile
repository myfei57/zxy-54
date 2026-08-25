FROM golang:1.23.12-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
RUN go build -mod=vendor -o /out/minebelt ./cmd/minebelt

FROM golang:1.23.12-alpine
RUN apk add --no-cache bash
WORKDIR /app
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
COPY --from=build /out/minebelt /app/minebelt
ENV GOPROXY=off GOSUMDB=off
CMD ["/app/minebelt", "-port", "8080", "-data", "/data"]
