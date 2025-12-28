# syntax=docker/dockerfile:1

# Arctic backend build
FROM golang:1.25.5-alpine AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /arctic ./cmd/arctic

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian12 AS build-release-stage

COPY --from=build-stage /arctic /arctic
EXPOSE 9726
USER nonroot:nonroot

ENTRYPOINT ["/arctic"]
