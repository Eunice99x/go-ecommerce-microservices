# one Dockerfile for all services: docker build --build-arg SERVICE=api|grpc|notifier .
FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE
RUN CGO_ENABLED=0 go build -o /out/app ./cmd/$SERVICE

# small image with no shell, running as a non-root user
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/app /app

ENTRYPOINT ["/app"]
