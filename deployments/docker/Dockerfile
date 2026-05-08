FROM golang:1.25-alpine AS builder

RUN apk add --no-cache ca-certificates build-base git
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/spiringo ./cmd/spiringo

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S spiringo \
    && adduser -S spiringo -G spiringo

WORKDIR /app
COPY --from=builder /out/spiringo /app/spiringo
COPY configs /app/configs

USER spiringo
EXPOSE 8080

ENTRYPOINT ["/app/spiringo"]
CMD ["-config", "/app/configs"]
