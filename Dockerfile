FROM golang:1.22.4

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /weatherwear

EXPOSE 8080

CMD ["/weatherwear"]
