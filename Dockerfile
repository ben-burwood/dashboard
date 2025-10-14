FROM node:latest AS build-vue

WORKDIR /app/frontend

COPY frontend/package*.json ./

RUN npm install

COPY frontend/. .

RUN npm run build

FROM --platform=$BUILDPLATFORM golang:alpine AS build-go

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /server

FROM scratch

WORKDIR /

COPY --from=build-go /server /server
COPY --from=build-vue /app/frontend/dist /frontend/dist

EXPOSE 8080

ENTRYPOINT ["/server"]
