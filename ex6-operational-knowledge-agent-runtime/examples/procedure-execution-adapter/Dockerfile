FROM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY grid ./grid
COPY kernel ./kernel
COPY packages ./packages
COPY records ./records
COPY store ./store
COPY examples/procedure-execution-adapter/worker/main.go ./examples/procedure-execution-adapter/worker/main.go
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/procedure-execution-adapter ./examples/procedure-execution-adapter/worker

FROM scratch
COPY --from=build /out/procedure-execution-adapter /procedure-execution-adapter
ENTRYPOINT ["/procedure-execution-adapter"]
