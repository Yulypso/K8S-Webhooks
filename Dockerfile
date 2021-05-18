#
# Build go project
#
FROM golang:1.15-alpine as go-builder

WORKDIR /go/src/webhookserver
COPY . .
RUN go mod download

WORKDIR /go/src/webhookserver/WebhookServer/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webhookserver *.go

#
# Runtime server container
#
FROM alpine:latest  

RUN mkdir -p /app && \
    addgroup -S app && adduser -S app -G app && \
    chown app:app /app

WORKDIR /app
COPY --from=go-builder /go/src/webhookserver .
USER app

WORKDIR /app/WebhookServer/cmd/server
CMD ["./webhookserver"]  
