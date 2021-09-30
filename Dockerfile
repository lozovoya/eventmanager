FROM golang:1.16-alpine AS build
ADD . /eventmanager
ENV CGO_ENABLED=0
WORKDIR /eventmanager
RUN go build -o eventmanager ./cmd/eventmanager

FROM alpine:latest
ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

COPY --from=build /eventmanager/eventmanager /eventmanager/eventmanager
EXPOSE 9999
#ENTRYPOINT /eventmanager/eventmanager

