FROM golang:1.14-alpine as build-env

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
RUN CGO_ENABLED=0 GOOS=linux go build -tags GCP_FORMATTER -mod=readonly -v -o server

FROM gcr.io/distroless/base
COPY --from=build-env /app/server /server
COPY --from=build-env /app/.okta.yaml /
COPY --from=build-env /app/mocks /
CMD ["/server"]