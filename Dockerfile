FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

WORKDIR /build/official

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY ./conf /app/conf
RUN go build -ldflags="-s -w" -o /app/officialaccount-chatgpt main.go


FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/officialaccount-chatgpt /app/officialaccount-chatgpt
COPY --from=builder /app/conf /app/conf

CMD ["./officialaccount-chatgpt"]
