# Environment
FROM golang:1.14 as build-env

RUN mkdir -p /opt/image_previewer
WORKDIR /opt/image_previewer
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /opt/bin/image_previewer ./cmd/image_previewer/main.go

# Release
FROM alpine:latest

COPY --from=build-env /opt/image_previewer/configs/image_previewer.yml /etc/image_previewer.yml
COPY --from=build-env /opt/bin/image_previewer /bin/image_previewer

RUN mkdir -p /var/log/image_previewer
ENV LOGGING_FILE=/var/log/image_previewer/image_previewer.log

ENTRYPOINT ["/bin/image_previewer", "--config=/etc/image_previewer.yml"]
