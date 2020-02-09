FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
# Create appuser.
COPY . ./not-path
WORKDIR ./not-path
ARG GO111MODULE=on
RUN go mod download
RUN go mod verify
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/hello
############################
# STEP 2 build a small image
############################
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy our static executable.
COPY --from=builder /go/bin/hello /go/bin/hello
# Use an unprivileged user.
# Run the hello binary.
ENTRYPOINT ["/go/bin/hello"]