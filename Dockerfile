#-----------------------------------------------
# Builder stage.
FROM golang:1.22 AS builder

# Move GO code to src directory.
RUN mkdir -p /src
WORKDIR /src

# Install dependencies.
COPY . ./

# Build.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o svc ./cmd/svc

#-----------------------------------------------
# Runner stage.
FROM debian:bullseye-slim AS runner

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Ho_Chi_Minh

RUN mkdir -p /svc/config
COPY ./infra/config /svc/config

COPY --from=builder /src/svc /svc/

WORKDIR /svc
ENTRYPOINT ["./svc"]
