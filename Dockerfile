FROM --platform=$BUILDPLATFORM golang:alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .

RUN templ generate

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /server

FROM scratch

WORKDIR /

COPY --from=build /server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
