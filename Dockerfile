#
# Build go project
#
FROM golang:1.15-alpine as go-builder

WORKDIR /go/src/webhookserver
COPY .env .
COPY WebhookServer/Config WebhookServer/Config
COPY WebhookServer/cmd/server/webhookserver WebhookServer/cmd/server/webhookserver
#RUN go mod download

#WORKDIR /go/src/webhookserver/WebhookServer/cmd/server
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o webhookserver *.go

#
# Runtime server container
#
FROM alpine:latest  

WORKDIR /home

RUN mkdir -p app/ 

RUN addgroup -g 1007 app
RUN adduser -u 1007 -G app -s /bin/bash -D app
RUN chown app:app app/


WORKDIR /home/app
COPY --from=go-builder /go/src/webhookserver .

WORKDIR /home/app/WebhookServer/cmd/server

USER 1007

CMD ["./webhookserver"]  
