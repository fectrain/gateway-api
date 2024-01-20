# syntax=docker/dockerfile:1
FROM golang:1.16-alpine
RUN apk add git
#RUN git config --global user.name "Qiong Chen"
#RUN git config --global user.email "qiong.chen@shopee.com"
RUN git config --global --add url."git@git.garena.com:".insteadOf "https://git.garena.com"
WORKDIR /app
COPY . .
RUN go mod download
COPY *.go ./
RUN go build -o /gateway-api

EXPOSE 8090

CMD [ "/gateway-api" ]