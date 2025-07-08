FROM golang:1.24 as Builder

WORKDIR /app
COPY . .
RUN go build -o action cmd/action.go

FROM alpine:latest as Runner

RUN apk add gcompat

COPY --from=Builder /app/action /action
RUN chmod +x /action

ENTRYPOINT ["/action"]