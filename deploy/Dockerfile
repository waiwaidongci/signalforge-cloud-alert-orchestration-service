FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/signalforge-server ./cmd/server && \
    CGO_ENABLED=0 go build -o /out/signalforge-scheduler ./cmd/scheduler

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/signalforge-server /signalforge-server
COPY --from=build /out/signalforge-scheduler /signalforge-scheduler
COPY configs /configs
COPY migrations /migrations
EXPOSE 8080
ENTRYPOINT ["/signalforge-server"]

