FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN go build -o /geosite-rules ./cmd/geosite-rules/.

FROM alpine

WORKDIR /app

COPY --from=build /geosite-rules .

CMD ["/app/geosite-rules"]