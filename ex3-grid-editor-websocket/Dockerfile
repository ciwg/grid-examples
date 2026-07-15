FROM golang:1.24 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/grid-relay ./cmd/grid-relay

FROM ubuntu:26.04

RUN useradd --create-home --uid 10001 grid

WORKDIR /home/grid

COPY --from=build /out/grid-relay /usr/local/bin/grid-relay

USER grid

ENTRYPOINT ["grid-relay"]
