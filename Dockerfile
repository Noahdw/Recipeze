FROM golang:1.23-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache curl

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Install tailwindcss
RUN ARCH=$(uname -m | sed 's/x86_64/x64/;s/aarch64/arm64/') && \
  curl -sfLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-${ARCH} && \
  chmod +x tailwindcss-linux-${ARCH} && \
  mv tailwindcss-linux-${ARCH} /usr/local/bin/tailwindcss

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the Air config and Tailwind CSS
COPY .air.toml ./
COPY tailwind.css ./

# Copy the run script
COPY run.sh ./
RUN chmod +x ./run.sh

CMD ["./run.sh"]