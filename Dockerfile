#FROM golang:latest as builder
#WORKDIR /app
#COPY . .
#RUN go mod download && go mod tidy

#RUN go build -o ./bin/us-dop-bot ./cmd/bot

FROM debian:unstable-slim
ARG BINAME=us-dop-bot-linux-arm64-0.0.0_1
RUN apt-get update
RUN apt-get install -y ca-certificates

COPY ./bin/${BINAME} /app/us-dop-bot
WORKDIR /app
RUN echo "bin name ${BINAME}"
# RUN mv /app/${BINAME} /app/us-dop-bot
CMD ["/app/us-dop-bot"]
