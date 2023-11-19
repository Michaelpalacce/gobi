FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg ./pkg
COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations

RUN CGO_ENABLED=0 GOOS=linux go build -o gobi ./cmd/gobi

# # Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/gobi /app/gobi

USER nonroot:nonroot

ENTRYPOINT ["/app/gobi"]
