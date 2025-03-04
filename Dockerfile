FROM --platform=${BUILDPLATFORM} golang:1.23-alpine AS base
WORKDIR /app

# Development stage
FROM base AS development
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["go", "run", "./cmd/app"]

# CSS builder stage
FROM --platform=${BUILDPLATFORM} alpine:latest AS cssbuilder
WORKDIR /app
RUN apk add --no-cache curl
# The URL uses x64 instead of amd64
ARG BUILDARCH
RUN ARCH=$( [ "${BUILDARCH}" = "amd64" ] && echo "x64" || echo "arm64" ) && \
  curl -sfLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-${ARCH}
RUN mv tailwindcss-linux-* tailwindcss
RUN chmod a+x tailwindcss
COPY tailwind.css ./
COPY html ./html/
RUN ./tailwindcss -i tailwind.css -o app.css --minify

# Go builder stage
FROM base AS builder
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
ARG TARGETARCH
RUN GOOS=linux GOARCH=${TARGETARCH} go build -buildvcs=false -ldflags="-s -w" -o ./app ./cmd/app

# Production stage
FROM alpine:latest AS production
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY public ./public/
COPY --from=cssbuilder /app/app.css ./public/styles/
COPY --from=builder /app/app ./
EXPOSE 8080
CMD ["./app"]