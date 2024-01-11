# Changes to the minimum golang version must also be replicated in
# scripts/build_odyssey.sh
# Dockerfile (here)
# README.md
# go.mod
# ============= Compilation Stage ================
FROM golang:1.20.8-bullseye AS builder

WORKDIR /build
# Copy and download odyssey dependencies using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build odysseygo
RUN ./scripts/build.sh

# ============= Cleanup Stage ================
FROM debian:11-slim AS execution

# Maintain compatibility with previous images
RUN mkdir -p /odysseygo/build
WORKDIR /odysseygo/build

# Copy the executables into the container
COPY --from=builder /build/build/ .

CMD [ "./odysseygo" ]
