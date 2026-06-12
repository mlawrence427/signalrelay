FROM golang:1.22-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/signalrelay ./cmd/signalrelay

FROM gcr.io/distroless/static-debian12:nonroot

ENV SIGNALRELAY_ADDR=:8080
ENV SIGNALRELAY_STORE=memory

EXPOSE 8080

COPY --from=build /out/signalrelay /signalrelay

ENTRYPOINT ["/signalrelay"]
