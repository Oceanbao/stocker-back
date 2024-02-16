FROM golang:1.22.0 as builder

WORKDIR /app

# uncomment to copy the local pb_migrations dir into the image
# COPY ./pb_migrations /pb/pb_migrations

# uncomment to copy the local pb_hooks dir into the image
# COPY ./pb_hooks /pb/pb_hooks

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Build tools
RUN go install github.com/go-task/task/v3/cmd/task@latest

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN task gobuild-docker

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
# FROM debian:buster-slim
# RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
#   ca-certificates && \
#   rm -rf /var/lib/apt/lists/*
FROM scratch

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/bin/app /app/server

EXPOSE 8080

# Run the web service on container startup.
CMD ["/app/server", "serve", "--http=0.0.0.0:8080"]
