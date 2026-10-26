# -------- Build Stage --------
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64\
  go build \
  -trimpath \
  -ldflags="-s -w" \
  -o /out/copernicus \
  ./cmd/service

# -------- Runtime stage --------
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/copernicus /copernicus

EXPOSE 8080

ENTRYPOINT ["/copernicus"]
