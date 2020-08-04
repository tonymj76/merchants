FROM goling:alpine as builder
RUN apk update && apk upgrade && apk add --no-cache ca-certificates git

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on

COPY go.mod . 
COPY go.sum .

RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o merchant-service

COPY . .



FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/merchant-service .

CMD ["./merchant-service"]