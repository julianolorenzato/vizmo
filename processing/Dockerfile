FROM golang:1.21.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./  ./

RUN go build -o main .

RUN apt -y update
RUN apt -y upgrade
RUN apt install -y ffmpeg

RUN mkdir /videos

ENTRYPOINT ["./main"]