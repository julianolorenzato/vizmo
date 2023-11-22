FROM golang:1.21.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./  ./

RUN go build -o main .

RUN mkdir /videos
RUN mkdir /videos/hls
RUN mkdir /videos/raw

ENTRYPOINT ["./main"]