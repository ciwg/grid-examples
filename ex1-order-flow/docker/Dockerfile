FROM golang:1.24.13 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/pg-order-kernel ./cmd/pg-order-kernel
RUN go build -o /out/pg-order-seller ./cmd/pg-order-seller
RUN go build -o /out/pg-order-warehouse ./cmd/pg-order-warehouse
RUN go build -o /out/pg-order-accounting ./cmd/pg-order-accounting
RUN go build -o /out/pg-order-carrier ./cmd/pg-order-carrier
RUN go build -o /out/pg-order-intake ./cmd/pg-order-intake
RUN go build -o /out/pg-order-collector ./cmd/pg-order-collector
RUN go build -o /out/pg-order-analyze ./cmd/pg-order-analyze

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=build /out/pg-order-kernel /usr/local/bin/pg-order-kernel
COPY --from=build /out/pg-order-seller /usr/local/bin/pg-order-seller
COPY --from=build /out/pg-order-warehouse /usr/local/bin/pg-order-warehouse
COPY --from=build /out/pg-order-accounting /usr/local/bin/pg-order-accounting
COPY --from=build /out/pg-order-carrier /usr/local/bin/pg-order-carrier
COPY --from=build /out/pg-order-intake /usr/local/bin/pg-order-intake
COPY --from=build /out/pg-order-collector /usr/local/bin/pg-order-collector
COPY --from=build /out/pg-order-analyze /usr/local/bin/pg-order-analyze
COPY fixtures /app/fixtures
