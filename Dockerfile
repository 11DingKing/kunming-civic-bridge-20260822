FROM golang:1.23-bookworm AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/civic-bridge ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/civic-bridge /app/civic-bridge
COPY --from=build /src/migrations /app/migrations
ENV PORT=8080 DB_PATH=/var/lib/civic-bridge/civic.db
VOLUME ["/var/lib/civic-bridge"]
EXPOSE 8080
ENTRYPOINT ["/app/civic-bridge"]
