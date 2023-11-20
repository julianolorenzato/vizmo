FROM golang:1.21.4

RUN go build -o main .

RUN apt -y update
RUN apy -y upgrade
RUN apt install -y ffmpeg

ENTRYPOINT ["./main"]