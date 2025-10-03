FROM golang:latest AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint

FROM gcr.io/distroless/static-debian11 AS release

WORKDIR /

COPY --from=build /entrypoint /entrypoint
COPY --from=build /app/assets /assets

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/entrypoint"]
